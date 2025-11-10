package pw

import (
	"bytes"
	"encoding/binary"

	"golang.org/x/crypto/scrypt"
)

type SiteVariant string
type Kind string

const (
	PASSWORD SiteVariant = "password"
	LOGIN    SiteVariant = "login"
	ANSWER   SiteVariant = "answer"

	MAXIMUM Kind = "maximum"
	LONG    Kind = "long"
	MEDIUM  Kind = "medium"
	BASIC   Kind = "basic"
	SHORT   Kind = "short"
	PIN     Kind = "pin"
	NAME    Kind = "name"
	PHRASE  Kind = "phrase"
)

func toBytes(n int) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(n))
	return b
}

func scopeForVariant(siteVariant SiteVariant) []byte {
	switch siteVariant {
	case PASSWORD:
		return []byte("com.lyndir.masterpassword")
	case LOGIN:
		return []byte("com.lyndir.masterpassword.login")
	default:
		return []byte("com.lyndir.masterpassword.answer")
	}
}

func templateOf(kind Kind, seed byte) []byte {
	templates := map[Kind][][]byte{
		MAXIMUM: {[]byte("anoxxxxxxxxxxxxxxxxx"), []byte("axxxxxxxxxxxxxxxxxno")},
		LONG: {
			[]byte("CvcvnoCvcvCvcv"),
			[]byte("CvcvCvcvnoCvcv"),
			[]byte("CvcvCvcvCvcvno"),
			[]byte("CvccnoCvcvCvcv"),
			[]byte("CvccCvcvnoCvcv"),
			[]byte("CvccCvcvCvcvno"),
			[]byte("CvcvnoCvccCvcv"),
			[]byte("CvcvCvccnoCvcv"),
			[]byte("CvcvCvccCvcvno"),
			[]byte("CvcvnoCvcvCvcc"),
			[]byte("CvcvCvcvnoCvcc"),
			[]byte("CvcvCvcvCvccno"),
			[]byte("CvccnoCvccCvcv"),
			[]byte("CvccCvccnoCvcv"),
			[]byte("CvccCvccCvcvno"),
			[]byte("CvcvnoCvccCvcc"),
			[]byte("CvcvCvccnoCvcc"),
			[]byte("CvcvCvccCvccno"),
			[]byte("CvccnoCvcvCvcc"),
			[]byte("CvccCvcvnoCvcc"),
			[]byte("CvccCvcvCvccno"),
		},
		MEDIUM: {[]byte("CvcnoCvc"), []byte("CvcCvcno")},
		BASIC:  {[]byte("aaanaaan"), []byte("aannaaan"), []byte("aaannaaa")},
		SHORT:  {[]byte("Cvcn")},
		PIN:    {[]byte("nnnn")},
		NAME:   {[]byte("cvccvcvcv")},
		PHRASE: {[]byte("cvcc cvc cvccvcv cvc"), []byte("cvc cvccvcvcv cvcv"), []byte("cv cvccv cvc cvcvccv")},
	}
	template := templates[kind]

	return template[int(seed)%len(template)]
}

func charFromClass(class byte, seed int) byte {
	lookup := map[byte][]byte{
		'V': []byte("AEIOU"),
		'C': []byte("BCDFGHJKLMNPQRSTVWXYZ"),
		'v': []byte("aeiou"),
		'c': []byte("bcdfghjklmnpqrstvwxyz"),
		'A': []byte("AEIOUBCDFGHJKLMNPQRSTVWXYZ"),
		'a': []byte("AEIOUaeiouBCDFGHJKLMNPQRSTVWXYZbcdfghjklmnpqrstvwxyz"),
		'n': []byte("0123456789"),
		'o': []byte("@&%?,=[]_:-+*$#!'^~;()/."),
		'x': []byte("AEIOUaeiouBCDFGHJKLMNPQRSTVWXYZbcdfghjklmnpqrstvwxyz0123456789!@#$%^&*()"),
		' ': []byte(" "),
	}
	bs := lookup[class]

	return bs[seed%len(bs)]
}

func mainKey(fullName, mainPass string, siteVariant SiteVariant) ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})
	defer buffer.Reset()

	buffer.Write(scopeForVariant(siteVariant))
	buffer.Write(toBytes(len(fullName)))
	buffer.Write([]byte(fullName))

	return scrypt.Key([]byte(mainPass), buffer.Bytes(), 32768, 8, 2, 64)
}

func password(mainKey []byte, site, context string, counter int, siteVariant SiteVariant, kind Kind) (string, error) {
	buffer := bytes.NewBuffer([]byte{})
	defer buffer.Reset()

	buffer.Write(scopeForVariant(siteVariant))
	buffer.Write(toBytes(len(site)))
	buffer.Write([]byte(site))
	buffer.Write(toBytes(counter))

	if context != "" {
		buffer.Write(toBytes(len(context)))
		buffer.Write([]byte(context))
	}

	seed, err := hmacSha256(mainKey, buffer.Bytes())
	if err != nil {
		return "", err
	}

	var passBytes []byte
	for i, ch := range templateOf(kind, seed[0]) {
		passBytes = append(passBytes, charFromClass(ch, int(seed[i+1])))
	}

	return string(passBytes), nil
}

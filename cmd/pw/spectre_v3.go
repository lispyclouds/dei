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
	var templates [][]byte

	switch kind {
	case MAXIMUM:
		templates = [][]byte{
			[]byte("anoxxxxxxxxxxxxxxxxx"),
			[]byte("axxxxxxxxxxxxxxxxxno"),
		}
	case LONG:
		templates = [][]byte{
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
		}
	case MEDIUM:
		templates = [][]byte{[]byte("CvcnoCvc"), []byte("CvcCvcno")}
	case BASIC:
		templates = [][]byte{
			[]byte("aaanaaan"),
			[]byte("aannaaan"),
			[]byte("aaannaaa"),
		}
	case SHORT:
		templates = [][]byte{[]byte("Cvcn")}
	case PIN:
		templates = [][]byte{[]byte("nnnn")}
	case NAME:
		templates = [][]byte{[]byte("cvccvcvcv")}
	case PHRASE:
		templates = [][]byte{
			[]byte("cvcc cvc cvccvcv cvc"),
			[]byte("cvc cvccvcvcv cvcv"),
			[]byte("cv cvccv cvc cvcvccv"),
		}
	}

	return templates[int(seed)%len(templates)]
}

func charFromClass(class byte, seed int) byte {
	var lookup []byte

	switch class {
	case 'V':
		lookup = []byte("AEIOU")
	case 'C':
		lookup = []byte("BCDFGHJKLMNPQRSTVWXYZ")
	case 'v':
		lookup = []byte("aeiou")
	case 'c':
		lookup = []byte("bcdfghjklmnpqrstvwxyz")
	case 'A':
		lookup = []byte("AEIOUBCDFGHJKLMNPQRSTVWXYZ")
	case 'a':
		lookup = []byte("AEIOUaeiouBCDFGHJKLMNPQRSTVWXYZbcdfghjklmnpqrstvwxyz")
	case 'n':
		lookup = []byte("0123456789")
	case 'o':
		lookup = []byte("@&%?,=[]_:-+*$#!'^~;()/.")
	case 'x':
		lookup = []byte("AEIOUaeiouBCDFGHJKLMNPQRSTVWXYZbcdfghjklmnpqrstvwxyz0123456789!@#$%^&*()")
	case ' ':
		lookup = []byte(" ")
	}

	return lookup[seed%len(lookup)]
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

package pw

import (
	"bytes"
	"encoding/binary"

	"golang.org/x/crypto/scrypt"
)

type SiteVariant string
type TemplateClass string

const (
	PASSWORD SiteVariant = "password"
	LOGIN    SiteVariant = "login"
	ANSWER   SiteVariant = "answer"

	MAXIMUM TemplateClass = "maximum"
	LONG    TemplateClass = "long"
	MEDIUM  TemplateClass = "medium"
	BASIC   TemplateClass = "basic"
	SHORT   TemplateClass = "short"
	PIN     TemplateClass = "pin"
	NAME    TemplateClass = "name"
	PHRASE  TemplateClass = "phrase"
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

func templateOf(class TemplateClass, seed byte) []byte {
	var templateSet [][]byte

	switch class {
	case MAXIMUM:
		templateSet = [][]byte{
			[]byte("anoxxxxxxxxxxxxxxxxx"),
			[]byte("axxxxxxxxxxxxxxxxxno"),
		}
	case LONG:
		templateSet = [][]byte{
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
		templateSet = [][]byte{[]byte("CvcnoCvc"), []byte("CvcCvcno")}
	case BASIC:
		templateSet = [][]byte{
			[]byte("aaanaaan"),
			[]byte("aannaaan"),
			[]byte("aaannaaa"),
		}
	case SHORT:
		templateSet = [][]byte{[]byte("Cvcn")}
	case PIN:
		templateSet = [][]byte{[]byte("nnnn")}
	case NAME:
		templateSet = [][]byte{[]byte("cvccvcvcv")}
	case PHRASE:
		templateSet = [][]byte{
			[]byte("cvcc cvc cvccvcv cvc"),
			[]byte("cvc cvccvcvcv cvcv"),
			[]byte("cv cvccv cvc cvcvccv"),
		}
	}

	return templateSet[int(seed)%len(templateSet)]
}

func charFromClass(class byte, seed int) byte {
	var charSet []byte

	switch class {
	case 'V':
		charSet = []byte("AEIOU")
	case 'C':
		charSet = []byte("BCDFGHJKLMNPQRSTVWXYZ")
	case 'v':
		charSet = []byte("aeiou")
	case 'c':
		charSet = []byte("bcdfghjklmnpqrstvwxyz")
	case 'A':
		charSet = []byte("AEIOUBCDFGHJKLMNPQRSTVWXYZ")
	case 'a':
		charSet = []byte("AEIOUaeiouBCDFGHJKLMNPQRSTVWXYZbcdfghjklmnpqrstvwxyz")
	case 'n':
		charSet = []byte("0123456789")
	case 'o':
		charSet = []byte("@&%?,=[]_:-+*$#!'^~;()/.")
	case 'x':
		charSet = []byte("AEIOUaeiouBCDFGHJKLMNPQRSTVWXYZbcdfghjklmnpqrstvwxyz0123456789!@#$%^&*()")
	case ' ':
		charSet = []byte(" ")
	}

	return charSet[seed%len(charSet)]
}

func mainKey(fullName, mainPass string, siteVariant SiteVariant) ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})
	defer buffer.Reset()

	buffer.Write(scopeForVariant(siteVariant))
	buffer.Write(toBytes(len(fullName)))
	buffer.Write([]byte(fullName))

	return scrypt.Key([]byte(mainPass), buffer.Bytes(), 32768, 8, 2, 64)
}

func password(mainKey []byte, site, context string, counter int, siteVariant SiteVariant, class TemplateClass) (string, error) {
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
	for i, ch := range templateOf(class, seed[0]) {
		passBytes = append(passBytes, charFromClass(ch, int(seed[i+1])))
	}

	return string(passBytes), nil
}

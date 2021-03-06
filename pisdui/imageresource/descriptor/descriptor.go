package descriptor

import (
	"os"

	"github.com/pisdhooy/fmtbytes"
)

type OsKeyBlock interface {
	getOsKeyBlockID() string
}

/*Descriptor block in an image resource*/
type Descriptor struct {
	Version       uint32
	UnicodeString string
	ClassID       string
	ItemCount     uint32
	Items         []*descriptorItem
}

type descriptorItem struct {
	key        string
	osTypeKey  string
	osKeyBlock OsKeyBlock
}

func (descriptor *Descriptor) GetTypeID() int {
	return 1088
}

func (Descriptor Descriptor) getOsKeyBlockID() string {
	return "objc"
}

/*NewDescriptor creates a new descriptor struct*/
func NewDescriptor() *Descriptor {
	return &Descriptor{}
}

/*Parse parses data from a descriptor block in a PSD file into a premade descriptor*/
func (descriptor *Descriptor) Parse(file *os.File) {

	descriptor.UnicodeString = fmtbytes.ParseUnicodeString(file)

	classIDLength := fmtbytes.ReadBytesLong(file)

	if classIDLength == 0 {
		descriptor.ClassID = fmtbytes.ReadBytesString(file, 4)
	} else {
		descriptor.ClassID = fmtbytes.ReadBytesString(file, int(classIDLength))
	}
	descriptor.ItemCount = fmtbytes.ReadBytesLong(file)

	var i uint32
	for i = 0; i < descriptor.ItemCount; i++ {
		descriptor.parseDescriptorItem(file)
	}

}

func (descriptor *Descriptor) parseDescriptorItem(file *os.File) {
	descriptorItem := new(descriptorItem)
	length := fmtbytes.ReadBytesLong(file)
	if length == 0 {
		descriptorItem.key = fmtbytes.ReadBytesString(file, 4)
	} else {
		descriptorItem.key = fmtbytes.ReadBytesString(file, int(length))
	}

	descriptorItem.osTypeKey = fmtbytes.ReadBytesString(file, 4)
	descriptorItem.osKeyBlock = parseOsKeyType(file, descriptorItem.osTypeKey)
	descriptor.Items = append(descriptor.Items, descriptorItem)
}

func parseOsKeyType(file *os.File, osKeyID string) OsKeyBlock {
	var r OsKeyBlock
	switch osKeyID {
	case "obj ":
		referenceObject := NewReference()
		referenceObject.Parse(file)
		break
	case "Objc":
		descriptorObject := NewDescriptor()
		descriptorObject.Parse(file)
		r = descriptorObject
		break
	case "VlLs":
		listObject := NewList()
		listObject.Parse(file)
		r = listObject
		break
	case "doub":
		doubleObject := NewDouble()
		err := doubleObject.Parse(file)
		if err != nil {
			//TODO: return this error properly
			panic(err)
		}
		r = doubleObject
		break
	case "UntF":
		unitFloatObject := NewUnitFloat()
		err := unitFloatObject.Parse(file)
		if err != nil {
			//TODO: return this error properly
			panic(err)
		}
		r = unitFloatObject
		break
	case "TEXT":
		textObject := NewText()
		textObject.Parse(file)
		r = textObject
		break
	case "enum":
		enumObject := NewEnum()
		enumObject.Parse(file)
		r = enumObject
		break
	case "long":
		integerObject := NewInteger()
		integerObject.Parse(file)
		r = integerObject
		break
	case "comp":
		break
	case "bool":
		boolObject := NewBool()
		boolObject.Parse(file)
		r = boolObject
		break
	case "GlbO":
		break
	case "type": //type and GlbC are both of type class
		fallthrough
	case "GlbC":
		break
	case "alis":
		break
	case "tdta":
		break
	}
	return r
}

//List is defined here to prevent cyclic imports to types.go
type List struct {
	NumItems uint32
	Items    []OsKeyBlock
}

func (list List) getOsKeyBlockID() string {
	return "VlLs"
}

func NewList() *List {
	return new(List)
}

func (list *List) Parse(file *os.File) {
	list.NumItems = fmtbytes.ReadBytesLong(file)
	for i := 0; i < int(list.NumItems); i++ {
		listItemType := fmtbytes.ReadBytesString(file, 4)
		list.Items = append(list.Items, parseOsKeyType(file, listItemType))
	}
}

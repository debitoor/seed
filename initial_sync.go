package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type NamespaceDoc struct {
	Name    string            `bson:name`
	Options CollectionOptions `bson:options`
}
type CollectionOptions struct {
	Name        string `bson:name`
	Capped      bool   `bson:capped`
	Size        int    `bson:size`
	Max         int    `bson:max`
	AutoIndexId bool   `bson:autoIndexId`
}

func (o CollectionOptions) String() string {
	return fmt.Sprintf("Capped: %b, Size: %d, Max: %d", o.Capped, o.Size, o.Max)
}

func CollectionInfo(c *mgo.Collection) (*mgo.CollectionInfo, error) {
	doc := NamespaceDoc{}

	if err := c.Database.C("system.namespaces").Find(bson.M{"name": c.FullName}).One(&doc); err != nil {
		return &mgo.CollectionInfo{}, err
	}

	return &mgo.CollectionInfo{
		Capped:         doc.Options.Capped,
		MaxBytes:       doc.Options.Size,
		MaxDocs:        doc.Options.Max,
		DisableIdIndex: !doc.Options.AutoIndexId}, nil
}

// find the size of a document in bytes
func docSize(ops interface{}) (int, error) {
	b, err := bson.Marshal(ops)
	if err != nil {
		return 0, err
	}
	return len(b), nil
}

package datadog

import (
"fmt"
"github.com/ooyala/go-dogstatsd"
}


type DataDogConf struct {
	Host     string
	Port     string
	Nameapce string
}

type Datadog struct {
	client *dogstatsd.Client
}

var datadogObj *Datadog

func GetInstance() *Datadog {
	if datadogObj == nil {
		//create data dog client connection
		c, err := dogstatsd.New("127.0.0.1:8125")
		if err != nil {
			fmt.Println("Could not connect to Datadog Agent", err)
		}
		//set NameSpace
		c.Namespace = "txRapid"

		//set tags
		var rTags []string
		rTags = append(rTags, "Rapid")
		c.Tags = rTags 


		datadogObj = new(Datadog)
		datadogObj.client = c

	}
	return client
}
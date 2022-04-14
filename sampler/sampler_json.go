package sampler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type SamplerJSON struct {
	Meta       SamplerMeta
	Samples_1D SampleData_1D
}

// Sampler Name
type SamplerMeta struct {
	Name      string `json:"name,omitempty"`
	SamplerID int    `json:"samplerID,omitempty"`
}

type SampleData_1D struct {
	Domain []float32 `json:"domain,omitempty"`
	Values []float32 `json:"values,omitempty"`
	Count  int       `json:"samples,count"`
}

func (strct *SamplerJSON) ExportJSON(filename string) {

	content, _ := strct.MarshalJSON()
	err := ioutil.WriteFile(filename, content, 0644)
	if err != nil {
		fmt.Printf("Unable to write JSON %s\n", filename)
	}
}

func (strct *SamplerMeta) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false

	buf.WriteString("\"name\": ")
	if tmp, err := json.Marshal(strct.Name); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
		comma = true
	}

	buf.WriteString("\"samplerID\": ")
	if tmp, err := json.Marshal(strct.SamplerID); err != nil {
		return nil, err
	} else {
		if comma {
			buf.WriteString(",")
		}
		buf.Write(tmp)
		comma = true
	}

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

func (strct *SamplerMeta) UnmarshalJSON(b []byte) error {
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "name":
			if err := json.Unmarshal([]byte(v), &strct.Name); err != nil {
				return err
			}

		case "samplerID":
			if err := json.Unmarshal([]byte(v), &strct.SamplerID); err != nil {
				return err
			}
		}
	}
	return nil
}

//------------------------------------------------------------------------------
//  SampleData_1D

func (strct *SampleData_1D) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false

	buf.WriteString("\"count\": ")
	if tmp, err := json.Marshal(strct.Count); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
		comma = true
	}

	buf.WriteString("\"domain\": ")
	if tmp, err := json.Marshal(strct.Domain); err != nil {
		return nil, err
	} else {
		if comma {
			buf.WriteString(",")
		}
		buf.Write(tmp)
		comma = true
	}

	buf.WriteString("\"values\": ")
	if tmp, err := json.Marshal(strct.Values); err != nil {
		return nil, err
	} else {
		if comma {
			buf.WriteString(",")
		}
		buf.Write(tmp)
		comma = true
	}

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

func (strct *SampleData_1D) UnmarshalJSON(b []byte) error {
	var jsonMap map[string]json.RawMessage
	countRecieved := false
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "count":
			if err := json.Unmarshal([]byte(v), &strct.Count); err != nil {
				return err
			}

		case "domain":
			if err := json.Unmarshal([]byte(v), &strct.Domain); err != nil {
				return err
			}
		case "values":
			if err := json.Unmarshal([]byte(v), &strct.Values); err != nil {
				return err
			}
		}
	}

	if !countRecieved {
		return fmt.Errorf("SampleData_1D.Count property not recieved")
	}
	return nil
}

//---------------------------SamplesJSON -----------------------------

//------------------------------------------------------------------------------
//  SampleData_1D

func (strct *SamplerJSON) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false

	buf.WriteString("\"meta\": ")
	if tmp, err := json.Marshal(strct.Meta); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
		comma = true
	}

	buf.WriteString("\"samples_1D\": ")
	if tmp, err := json.Marshal(strct.Samples_1D); err != nil {
		return nil, err
	} else {
		if comma {
			buf.WriteString(",")
		}
		buf.Write(tmp)
		comma = true
	}

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

func (strct *SamplerJSON) UnmarshalJSON(b []byte) error {
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "meta":
			if err := json.Unmarshal([]byte(v), &strct.Meta); err != nil {
				return err
			}

		case "samples_1D":
			if err := json.Unmarshal([]byte(v), &strct.Samples_1D); err != nil {
				return err
			}
		}
	}

	return nil
}

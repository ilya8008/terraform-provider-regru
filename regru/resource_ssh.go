package main

import (
    "github.com/hashicorp/terraform/helper/schema"
	"net/http"
	"encoding/json"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"time"
)

// Type definitions for json 

type SSHKeyInit struct {
	Name  string `json:"name"`
	PublicKey  string `json:"public_key"`
}

type SshKeys []struct {
	ID  int        `json:"id"`
	Fingerprint string `json:"fingerprint"`
	Name string  `json:"name"`
	PublicKey  string `json:"public_key"`
}

type SSHReply struct {
	SshKeys `json:"ssh_keys"`
}

var g SSHReply

func resourceSSH() *schema.Resource {
        return &schema.Resource{
                Create: resourceSSHCreate,
                Read:   resourceSSHRead,
                Update: resourceSSHUpdate,
                Delete: resourceSSHDelete,

                Schema: map[string]*schema.Schema{
                        "name": &schema.Schema{
                                Type:     schema.TypeString,
                                Required: true,
                        },
                        "public_key": &schema.Schema{
                                Type:     schema.TypeString,
                                Required: true,
                        },
                        "token": &schema.Schema{
                                Type:     schema.TypeString,
                                Required: true,
                        },
                },
        }
}

func resourceSSHCreate(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	public_key := d.Get("public_key").(string)
	token := d.Get("token").(string)
	data := SSHKeyInit{name,public_key}
	payloadBytes, err := json.Marshal(data)
	if err != nil {

	}
	body := bytes.NewReader(payloadBytes)
	fmt.Println(string(payloadBytes))
	req, err := http.NewRequest("POST", "https://api.cloudvps.reg.ru/v1/account/keys", body)
	if err != nil {

	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
	}
	defer resp.Body.Close()
	d.SetId(getsshid(token, d))
	return nil
}

func resourceSSHRead(d *schema.ResourceData, m interface{}) error {
	token := d.Get("token").(string)
	res := getsshstatus(token,d)
	if res == "" {
		log.Printf("[WARN] No Key found: %s", d.Id())
		d.SetId("")
	}
	return nil
}

func resourceSSHUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceSSHRead(d, m)
}

func resourceSSHDelete(d *schema.ResourceData, m interface{}) error {
	token := d.Get("token").(string)
	url := "https://api.cloudvps.reg.ru/v1/account/keys/" + d.Id()
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		// handle err
	}
	req.Header.Set("Authorization", "Bearer " + d.Get("token").(string))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	defer resp.Body.Close()
	for ok := true; ok; ok = getsshstatus(token,d) != "" {
		time.Sleep(100 * 1000);
	    }
	d.SetId("")
	return nil
}

func getsshid(token string, d *schema.ResourceData) string {
	
	var id string
	req, err := http.NewRequest("GET", "https://api.cloudvps.reg.ru/v1/account/keys", nil)
	if err != nil {
		// handle err
	}
	req.Header.Set("Authorization", "Bearer " + token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &g)
	for index, value := range g.SshKeys {
		fmt.Println(string(index))
		if value.Name == d.Get("name").(string) {
			id = fmt.Sprintf("%v",value.ID)
		}
	}
	return(id)
}

func getsshstatus(token string, d *schema.ResourceData) string {
	
	var status string
	req, err := http.NewRequest("GET", "https://api.cloudvps.reg.ru/v1/account/keys", nil)
	if err != nil {
		// handle err
	}
	req.Header.Set("Authorization", "Bearer " + token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &g)
	for index, value := range g.SshKeys {
		fmt.Println(string(index))
		if fmt.Sprintf("%v",value.ID) == d.Id() {
			status = fmt.Sprintf("%v",value.ID)
		}
	}
	return(status)
}

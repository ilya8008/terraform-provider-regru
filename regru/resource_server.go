package regru

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Type definitions for json

type Payload struct {
	Name  string `json:"name"`
	Size  string `json:"size"`
	Image string `json:"image"`
}

type Size struct {
	Archived   bool   `json:"archived"`
	Disk       int    `json:"disk"`
	ID         int    `json:"id"`
	Memory     int    `json:"memory"`
	Name       string `json:"name"`
	Price      string `json:"price"`
	PriceMonth string `json:"price_month"`
	Slug       string `json:"slug"`
	Vcpus      int    `json:"vcpus"`
	Weight     int    `json:"weight"`
}

type Image struct {
	CreatedAt     string `json:"created_at"`
	Distribution  string `json:"distribution"`
	ID            int    `json:"id"`
	MinDiskSize   string `json:"min_disk_size"`
	Name          string `json:"name"`
	Private       bool   `json:"private"`
	RegionSlug    string `json:"region_slug"`
	SizeGigabytes string `json:"size_gigabytes"`
	Slug          string `json:"slug"`
	Type          string `json:"type"`
}

type Reglets []struct {
	ArchivedAt     string `json:"archived_at"`
	BackupsEnabled bool   `json:"backups_enabled"`
	CreatedAt      string `json:"created_at"`
	Disk           int    `json:"disk"`
	Hostname       string `json:"hostname"`
	ID             int    `json:"id"`
	Image          `json:"image"`
	ImageID        int    `json:"image_id"`
	IP             string `json:"ip"`
	Ipv6           string `json:"ipv6"`
	LinkToken      string `json:"link_token"`
	Locked         bool   `json:"locked"`
	Memory         int    `json:"memory"`
	Name           string `json:"name"`
	Ptr            string `json:"ptr"`
	RegionSlug     string `json:"region_slug"`
	ServiceID      int    `json:"service_id"`
	Size           `json:"size"`
	SizeSlug       string `json:"size_slug"`
	Status         string `json:"status"`
	SubStatus      string `json:"sub_status"`
	Vcpus          int    `json:"vcpus"`
}

type Actions []struct {
	CompletedAt  string  `json:"completed_at"`
	CreatedAt    string  `json:"created_at"`
	ID           float64 `json:"id"`
	RegionSlug   string  `json:"region_slug"`
	ResourceID   float64 `json:"resource_id"`
	ResourceType string  `json:"resource_type"`
	Status       string  `json:"status"`
	Type         string  `json:"type"`
}

type Links struct {
	Actions `json:"actions"`
}
type Reply struct {
	Links   `json:"links"`
	Reglets `json:"reglets"`
}

var f Reply

func resourceServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerCreate,
		Read:   resourceServerRead,
		Update: resourceServerUpdate,
		Delete: resourceServerDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"size": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"image": &schema.Schema{
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

func resourceServerCreate(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	size := d.Get("size").(string)
	image := d.Get("image").(string)
	token := d.Get("token").(string)
	data := Payload{name, size, image}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		// handle err
	}
	body := bytes.NewReader(payloadBytes)
	fmt.Println(string(payloadBytes))
	req, err := http.NewRequest("POST", "https://api.cloudvps.reg.ru/v1/reglets", body)
	if err != nil {
		// handle err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	defer resp.Body.Close()
	d.SetId(getserverid(token, d))
	for ok := true; ok; ok = getserverstatus(token, d) != "active" {
		time.Sleep(100 * 1000)
	}
	return nil
}

func resourceServerRead(d *schema.ResourceData, m interface{}) error {
	token := d.Get("token").(string)
	res := getserverstatus(token, d)
	if res == "" {
		log.Printf("[WARN] No Server found: %s", d.Id())
		d.SetId("")
	}
	return nil
}

func resourceServerUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceServerRead(d, m)
}

func resourceServerDelete(d *schema.ResourceData, m interface{}) error {
	token := d.Get("token").(string)
	url := "https://api.cloudvps.reg.ru/v1/reglets/" + d.Id()
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		// handle err
	}
	req.Header.Set("Authorization", "Bearer "+d.Get("token").(string))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	defer resp.Body.Close()
	for ok := true; ok; ok = getserverstatus(token, d) != "" {
		time.Sleep(100 * 1000)
	}
	d.SetId("")
	return nil
}

func getserverid(token string, d *schema.ResourceData) string {

	var id string
	req, err := http.NewRequest("GET", "https://api.cloudvps.reg.ru/v1/reglets", nil)
	if err != nil {
		// handle err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &f)
	for index, value := range f.Reglets {
		fmt.Println(string(index))
		if value.Name == d.Get("name").(string) {
			id = fmt.Sprintf("%v", value.ID)
		}
	}
	return (id)
}

func getserverstatus(token string, d *schema.ResourceData) string {

	var status string
	req, err := http.NewRequest("GET", "https://api.cloudvps.reg.ru/v1/reglets", nil)
	if err != nil {
		// handle err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &f)
	for index, value := range f.Reglets {
		fmt.Println(string(index))
		if fmt.Sprintf("%v", value.ID) == d.Id() {
			status = fmt.Sprintf("%v", value.Status)
		}
	}
	return (status)
}

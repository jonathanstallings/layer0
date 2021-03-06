package main

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/quintilesims/layer0/cli/client"
)

func resourceLayer0Environment() *schema.Resource {
	return &schema.Resource{
		Create: resourceLayer0EnvironmentCreate,
		Read:   resourceLayer0EnvironmentRead,
		Update: resourceLayer0EnvironmentUpdate,
		Delete: resourceLayer0EnvironmentDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"size": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "m3.medium",
				ForceNew: true,
			},
			"min_count": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"user_data": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"cluster_count": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_group_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceLayer0EnvironmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.APIClient)

	name := d.Get("name").(string)
	size := d.Get("size").(string)
	minCount := d.Get("min_count").(int)
	userData := d.Get("user_data").(string)

	environment, err := client.CreateEnvironment(name, size, minCount, []byte(userData))
	if err != nil {
		return err
	}

	d.SetId(environment.EnvironmentID)
	return resourceLayer0EnvironmentRead(d, meta)
}

func resourceLayer0EnvironmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.APIClient)
	environmentID := d.Id()

	environment, err := client.GetEnvironment(environmentID)
	if err != nil {
		if strings.Contains(err.Error(), "No environment found") {
			d.SetId("")
			log.Printf("[WARN] Error Reading Environment (%s), environment does not exist", environmentID)
			return nil
		}

		return err
	}

	d.Set("name", environment.EnvironmentName)
	d.Set("size", environment.InstanceSize)
	d.Set("cluster_count", environment.ClusterCount)
	d.Set("security_group_id", environment.SecurityGroupID)

	return nil
}

func resourceLayer0EnvironmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.APIClient)
	environmentID := d.Id()

	if d.HasChange("min_count") {
		minCount := d.Get("min_count").(int)

		if _, err := client.UpdateEnvironment(environmentID, minCount); err != nil {
			return err
		}
	}

	return resourceLayer0EnvironmentRead(d, meta)
}

func resourceLayer0EnvironmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.APIClient)
	environmentID := d.Id()

	jobID, err := client.DeleteEnvironment(environmentID)
	if err != nil {
		if strings.Contains(err.Error(), "No environment found") {
			return nil
		}

		return err
	}

	if err := client.WaitForJob(jobID, defaultTimeout); err != nil {
		return err
	}

	return nil
}

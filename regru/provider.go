package regru

import (
        "github.com/hashicorp/terraform/helper/schema"
)

func Provider() *schema.Provider {
        return &schema.Provider{
                ResourcesMap: map[string]*schema.Resource{
                        "regru_server": resourceServer(),
                        "regru_ssh": resourceSSH(),
                },
        }
}

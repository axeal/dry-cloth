package drycloth

import (
	"context"
	"log"
	"time"

	"github.com/digitalocean/godo"
)

func contains(arr []string, str string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}

func parseDate(str string) (time.Time, error) {
	layout := "2006-01-02T15:04:05Z"
	t, err := time.Parse(layout, str)
	return t, err
}

func Run(ctx context.Context, accessToken string, preserveTag string, days int, dryRun bool) error {
	client := godo.NewFromToken(accessToken)
	droplets, err := DropletList(ctx, client)
	if err != nil {
		log.Printf("Error listing droplets: %s", err)
		return err
	}

	pruneBefore := time.Now().Add(-time.Duration(days) * 24 * time.Hour)
	log.Printf("Pruning droplets created before %s", pruneBefore.Format("2006-01-02T15:04:05Z"))
	deleted := 0

	for _, droplet := range droplets {
		t, err := parseDate(droplet.Created)
		if err != nil {
			log.Printf("Error parsing created string (%s) for droplet: %s", droplet.Created, droplet.Name)
			continue
		}
		if t.Before(pruneBefore) {
			if preserveTag != "" && contains(droplet.Tags, preserveTag) {
				log.Printf("Skipping deletion of droplet %s, preservation tag %s is present", droplet.Name, preserveTag)
				continue
			} else {
				if dryRun {
					log.Printf("Droplet %s (created %s) would be deleted", droplet.Name, droplet.Created)
				} else {
					log.Printf("Deleting droplet %s (created %s)", droplet.Name, droplet.Created)
					_, err := client.Droplets.Delete(ctx, droplet.ID)
					if err != nil {
						log.Printf("Error deleting droplet %s : %s", droplet.Name, err)
					}
				}
			}
		}
	}

	return nil
}

// https://github.com/digitalocean/godo#pagination
func DropletList(ctx context.Context, client *godo.Client) ([]godo.Droplet, error) {
	// create a list to hold our droplets
	list := []godo.Droplet{}

	// create options. initially, these will be blank
	opt := &godo.ListOptions{}
	for {
		droplets, resp, err := client.Droplets.List(ctx, opt)
		if err != nil {
			return nil, err
		}

		// append the current page's droplets to our list
		list = append(list, droplets...)

		// if we are at the last page, break out the for loop
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, err
		}

		// set the page we want for the next request
		opt.Page = page + 1
	}

	return list, nil
}

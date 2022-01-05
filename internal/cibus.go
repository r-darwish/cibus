package internal

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
	"log"
	"time"
)

func AddAllFriends(username, password string) error {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.Navigate("https://www.mysodexo.co.il"),
		chromedp.SetValue("#txtUsr", username, chromedp.ByID),
		chromedp.SetValue("#txtPas", password, chromedp.ByID),
		chromedp.Click("#btnLogin"),
		chromedp.WaitVisible("#ctl00_lnkRound1"),
	)
	if err != nil {
		return fmt.Errorf("failed logging in to Cibus: %w", err)
	}

	var friends []*cdp.Node
	err = chromedp.Run(ctx,
		chromedp.Navigate("https://www.mysodexo.co.il/new_my/new_my_friends.aspx"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			for {
				deadline, c := context.WithTimeout(ctx, time.Second)
				err := chromedp.Nodes(`//div[contains(@class, "all-friends")]/div[contains(@class, "friends-panel")]/span`, &friends).Do(deadline)
				c()
				if err != nil {
					return err
				}
				if len(friends) == 0 {
					return nil
				}

				friend := friends[0]

				err = dom.RequestChildNodes(friend.NodeID).WithDepth(-1).Do(ctx)
				if err != nil {
					log.Printf("Error clicking a user: %s", err)
					continue
				}

				time.Sleep(1 * time.Second)
				if len(friend.Children) == 0 {
					log.Printf("Unclickable node")
					continue
				}
				log.Printf("Adding friend")

				deadline, c = context.WithTimeout(ctx, time.Second)
				err = chromedp.Click([]cdp.NodeID{friend.Children[1].NodeID}, chromedp.ByNodeID).Do(deadline)
				c()
				if err != nil {
					log.Printf("Error clicking a user: %s", err)
					continue
				}

				err = chromedp.WaitNotVisible(".ui-loader").Do(ctx)
				if err != nil {
					return err
				}

			}
		}),
	)
	if err != nil {
		return fmt.Errorf("could not get friends: %w", err)
	}

	return nil
}

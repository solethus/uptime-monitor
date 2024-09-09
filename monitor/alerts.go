package monitor

import (
	"context"
	"database/sql"
	"encore.app/site"
	"encore.dev/pubsub"
	"errors"
)

// TransitionEvent describes a transition of a monitored site from up->down or from down->up.
type TransitionEvent struct {
	// Site is the monitored site in question.
	Site *site.Site `json:"site"`
	// Up specifies whether a site is now up or down (the new value).
	Up bool `json:"up"`
}

// TransitionTopic is a pubsub topic with the transition events for when a monitored site
// transitions from up->down or from down->up.
var TransitionTopic = pubsub.NewTopic[*TransitionEvent]("uptime-transition", pubsub.TopicConfig{
	DeliveryGuarantee: pubsub.AtLeastOnce,
})

// getPreviousMeasurement reports whether the given site was up or down in previous measurement
func getPreviousMeasurement(ctx context.Context, siteID int) (wasUp bool, err error) {
	err = db.QueryRow(ctx, `
		SELECT up FROM checks
		WHERE site_id = $1
		ORDER BY checked_at DESC
		LIMIT 1
   	`, siteID).Scan(&wasUp)
	if errors.Is(err, sql.ErrNoRows) {
		// There was no previous ping; treat this as if the sites was up before.
		return true, nil
	} else if err != nil {
		return false, err
	}
	return wasUp, nil
}

func pushOnTransition(ctx context.Context, site *site.Site, isUp bool) error {
	wasUp, err := getPreviousMeasurement(ctx, site.ID)
	if err != nil {
		return err
	}
	if isUp == wasUp {
		// Nothing to do.
		return nil
	}

	_, err = TransitionTopic.Publish(ctx, &TransitionEvent{
		Site: site,
		Up:   isUp,
	})

	return err
}

package app

import (
	"context"
	"github.com/asim/go-micro/v3/metadata"
	"github.com/asim/go-micro/v3/server"
	"micro-libs/utils/errors"
	"micro-libs/utils/log"
	"time"
)

func serverWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		if Opts.Dev {
			mt, _ := metadata.FromContext(ctx)
			log.Info("[RPC] Call [%s] from [%s] ...", req.Endpoint(), mt["Micro-From-Service"])
		}

		now := time.Now()
		err := fn(ctx, req, rsp)

		if Opts.Dev {
			if err != nil {
				if errors.IsStack() {
					log.Error("[RPC] Received [%s], uptime: %s, error: %+v", req.Endpoint(), time.Since(now), err)
				} else {
					log.Error("[RPC] Received [%s], uptime: %s, error: %s", req.Endpoint(), time.Since(now), err)
				}
			} else {
				log.Info("[RPC] Received [%s], uptime: %s", req.Endpoint(), time.Since(now))
			}
		}

		if err == nil {
			return nil
		}

		return errors.MicroError(Id(), err)
	}
}

func subscriberWrapper(fn server.SubscriberFunc) server.SubscriberFunc {
	return func(ctx context.Context, msg server.Message) error {
		if Opts.Dev {
			mt, _ := metadata.FromContext(ctx)
			log.Info("[TOPIC] Call [%s] from [%s] ...", msg.Topic(), mt["Micro-From-Service"])
		}

		now := time.Now()
		err := fn(ctx, msg)

		if Opts.Dev {
			if err != nil {
				if errors.IsStack() {
					log.Error("[TOPIC] Received [%s], uptime: %s, error: %+v", msg.Topic(), time.Since(now), err)
				} else {
					log.Error("[TOPIC] Received [%s], uptime: %s, error: %s", msg.Topic(), time.Since(now), err)
				}
			} else {
				log.Info("[TOPIC] Received [%s], uptime: %s", msg.Topic(), time.Since(now))
			}
		}

		if err == nil {
			return nil
		}

		return errors.MicroError(Id(), err)
	}
}

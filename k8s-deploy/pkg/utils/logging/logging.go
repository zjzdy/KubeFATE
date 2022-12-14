/*
 * Copyright 2019-2021 VMware, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
package logging

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// InitLog initials a log instance with specified config
func InitLog() {

	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			// TODO(JHC):
			// 1. output device should read from config file
			Out:        os.Stdout,
			NoColor:    viper.GetBool("log.nocolor"),
			TimeFormat: time.RFC3339,
		},
	).With().Caller().Stack().Logger()
	logLevel, err := zerolog.ParseLevel(viper.GetString("log.level"))
	if err != nil {
		log.Error().Err(err).
			Str("You need to choose one from here", fmt.Sprint(
				zerolog.TraceLevel,
				zerolog.InfoLevel,
				zerolog.DebugLevel,
				zerolog.WarnLevel,
				zerolog.ErrorLevel,
				zerolog.FatalLevel,
				zerolog.PanicLevel,
			)).
			Msg("Get log level configuration error")
	}
	log.Logger = log.Logger.Level(logLevel)
}

var GetGinLogger = func(c *gin.Context, out io.Writer, latency time.Duration) zerolog.Logger {
	logger := zerolog.New(out).With().
		Timestamp().
		Int("status", c.Writer.Status()).
		Str("method", c.Request.Method).
		Str("path", c.Request.URL.Path).
		Str("ip", c.ClientIP()).
		Dur("latency", latency).
		Str("user_agent", c.Request.UserAgent()).
		Logger().Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		NoColor:    viper.GetBool("log.nocolor"),
		TimeFormat: time.RFC3339,
	})
	return logger

}

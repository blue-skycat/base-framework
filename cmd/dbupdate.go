// Copyright 2018 cloudy 272685110@qq.com.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package cmd

import (
	"errors"
	"github.com/itcloudy/base-framework/pkg/conf"
	"github.com/itcloudy/base-framework/pkg/logs"
	"github.com/itcloudy/base-framework/pkg/repositories/common"
	"github.com/itcloudy/base-framework/pkg/services"
	"go.uber.org/zap"
	"time"

	"github.com/spf13/cobra"
)

var dbupdateCmd = &cobra.Command{
	Use:    "dbupdate",
	Short:  "update database, add tables and columns or modify columns",
	PreRun: loadConfig,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := conf.Config.DB
		conf.GetDBConnection(cfg.DbType, cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.Charset, "update")
		migrateService := services.MigrationService{}
		switch cfg.DbType {
		case "mysql":
		case "postgres":
			migrateService.IMigrationHistoryRepository = &common.MigrationHistoryRepository{DB: conf.SqlxDB}
			break
		default:
			panic(errors.New("un support sql type:" + cfg.DbType))

		}
		if err := migrateService.UpdateToOneVersion(); err != nil {
			logs.Logger.Fatal("update database failed", zap.String("db name", conf.Config.DB.Name), zap.Error(err))
			logs.Logger.Sync()

		} else {
			logs.Logger.Info("update database success", zap.String("db name", conf.Config.DB.Name), zap.Time("time", time.Now()))
			logs.Logger.Sync()

		}
	},
}

func init() {
	dbupdateCmd.Flags().StringVar(&conf.Config.DBUpdateToVersion, "to-version", "", "update database to version")
	dbupdateCmd.MarkFlagRequired("to-version")
}
package com.jiny.builder;

import com.alibaba.druid.pool.DruidDataSource;
import com.taosdata.jdbc.TSDBDriver;
import lombok.extern.slf4j.Slf4j;

import javax.sql.DataSource;
import java.util.Properties;

/**
 * @Auther: jiny
 * @CreateDate: 2022/7/5 14:03
 * @Description: 连接池
 */
@Slf4j
public class TaosDruidPool {

    public static DataSource DATASOURCE;

    static {
        DATASOURCE = getDataSource("taos-server", 50, "root", "taosdata", "h");
    }

    private static DataSource getDataSource(String host, int poolSize, String username, String password, String database) {
        String url = "jdbc:TAOS://" + host + ":6030/" + database;
        DruidDataSource dataSource = new DruidDataSource();
        // jdbc properties
        dataSource.setDriverClassName("com.taosdata.jdbc.TSDBDriver");
        dataSource.setUrl(url);
        dataSource.setUsername(username);
        dataSource.setPassword(password);
        Properties connProps = new Properties();
        connProps.setProperty(TSDBDriver.PROPERTY_KEY_BATCH_LOAD, "true");
        connProps.setProperty(TSDBDriver.PROPERTY_KEY_CHARSET, "utf-8");
        connProps.setProperty(TSDBDriver.PROPERTY_KEY_LOCALE, "en_US.UTF-8");
        connProps.setProperty(TSDBDriver.PROPERTY_KEY_TIME_ZONE, "UTC-8");
        //connProps.setProperty("debugFlag", "135");
        connProps.setProperty("maxSQLLength", "1048576");
        dataSource.setConnectProperties(connProps);
        // pool configurations
        dataSource.setInitialSize(poolSize);
        dataSource.setMinIdle(poolSize);
        dataSource.setMaxActive(poolSize);
        dataSource.setMaxWait(30000);
        dataSource.setValidationQuery("select server_status()");
        log.info("taos pool url {}", dataSource.getUrl());
        return dataSource;
    }
}

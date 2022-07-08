package com.jiny.service;

import com.jiny.api.TaosInterface;
import com.jiny.builder.TaosDruidPool;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import org.springframework.util.StringUtils;

import java.sql.Connection;
import java.sql.SQLException;
import java.sql.Statement;
import java.text.MessageFormat;

/**
 * @Auther: jiny
 * @CreateDate: 2022/7/5 10:36
 * @Description:
 */
@Service
@Slf4j
public class TaosImpl implements TaosInterface {


    @Override
    public Boolean updateSubTableTag(String subTable, String tagName, String tagValue) {
        if (!StringUtils.hasText(subTable) || !StringUtils.hasText(tagName) || !StringUtils.hasText(tagValue)) {
            return false;
        }
        try {
            String sql = MessageFormat.format("ALTER TABLE {0} SET TAG {1}={2};", subTable, tagName, tagValue);
            executeDDL(sql);
            return true;
        } catch (SQLException e) {
            log.error("throws error: code:{},msg:{}", e.getErrorCode(), e.getMessage());
            e.printStackTrace();
            return false;
        }
    }

    @Override
    public Boolean insertSubTable() {
        return null;
    }


    private void executeDDL(String ddlSql) throws SQLException {
        try (Connection connection = getConnection();
             Statement statement = connection.createStatement()) {
            statement.executeUpdate(ddlSql);
            //resetCache();
        }
    }

    private void resetCache() throws SQLException {
        /*try (Connection connection = getConnection();
             Statement statement = connection.createStatement()) {
            statement.executeUpdate("reset query cache;");
        }*/

    }


    /**
     * @Author: jiny
     * @CreateDate: 2022/7/7 9:58
     * @Description: 获取连接
     */
    private Connection getConnection() throws SQLException {
        return TaosDruidPool.DATASOURCE.getConnection();
    }



}

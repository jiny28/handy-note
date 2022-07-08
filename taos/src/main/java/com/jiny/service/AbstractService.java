package com.jiny.service;

import com.jiny.builder.TaosDruidPool;

import java.sql.*;

/**
 * @Auther: jiny
 * @CreateDate: 2022/7/8 17:11
 * @Description:
 */
public class AbstractService {


    public void executeDDL(String sql) throws SQLException {
        try (Connection connection = getConnection();
             Statement statement = connection.createStatement()) {
            statement.executeUpdate(sql);
        }
    }



    public Connection getConnection() throws SQLException {
        return TaosDruidPool.DATASOURCE.getConnection();
    }
}

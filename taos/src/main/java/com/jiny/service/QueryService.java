package com.jiny.service;

import com.jiny.api.QueryInterface;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import org.springframework.util.StringUtils;

import java.sql.*;
import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;

/**
 * @Auther: jiny
 * @CreateDate: 2022/7/11 16:17
 * @Description:
 */
@Service
@Slf4j
public class QueryService extends AbstractService implements QueryInterface {

    @Override
    public List<LinkedHashMap<String, String>> executeQuery(String sql) {
        if (!StringUtils.hasText(sql)) return null;
        log.debug("SQL >>> " + sql);
        List<LinkedHashMap<String, String>> result = new ArrayList<>();
        try (Connection connection = getConnection();
             Statement statement = connection.createStatement()) {
            ResultSet resultSet = statement.executeQuery(sql);
            ResultSetMetaData meta = resultSet.getMetaData();
            while (resultSet.next()) {
                LinkedHashMap<String, String> map = new LinkedHashMap<>();
                for (int i = 1; i <= meta.getColumnCount(); i++) {
                    String columnLabel = meta.getColumnLabel(i);
                    String value;
                    if (meta.getColumnType(i) == Types.TIMESTAMP) {
                        Timestamp timestamp = resultSet.getTimestamp(i);
                        value = timestamp.getTime() + "";
                    } else {
                        value = resultSet.getString(i);
                    }
                    map.put(columnLabel, value);
                }
                result.add(map);
            }
        } catch (SQLException e) {
            log.error("throws error: code:{},msg:{}", e.getErrorCode(), e.getMessage());
            e.printStackTrace();
        }
        return result;
    }

}

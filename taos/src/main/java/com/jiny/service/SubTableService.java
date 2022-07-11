package com.jiny.service;

import com.jiny.api.SubTableInterface;
import com.jiny.entity.SubTableMeta;
import com.jiny.entity.SubTableValue;
import com.jiny.utils.SqlSpeller;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import org.springframework.util.StringUtils;

import java.sql.SQLException;
import java.text.MessageFormat;
import java.util.List;

/**
 * @Auther: jiny
 * @CreateDate: 2022/7/11 11:33
 * @Description:
 */
@Service
@Slf4j
public class SubTableService extends AbstractService implements SubTableInterface {


    @Override
    public void create(SubTableMeta subTableMeta) {
        String sql = SqlSpeller.createTableUsingSuperTable(subTableMeta);
        log.debug("SQL >>> " + sql);
        try {
            executeDDL(sql);
        } catch (SQLException e) {
            log.error("throws error: code:{},msg:{}", e.getErrorCode(), e.getMessage());
            e.printStackTrace();
        }
    }

    @Override
    public int insert(List<SubTableValue> subTableValues) {
        String sql = SqlSpeller.insertMultiSubTableMultiValues(subTableValues);
        log.debug("SQL >>> " + sql);
        int affectRows = 0;
        try {
            affectRows = executeDDL(sql);
        } catch (SQLException e) {
            log.error("throws error: code:{},msg:{}", e.getErrorCode(), e.getMessage());
            e.printStackTrace();
        }
        return affectRows;
    }

    @Override
    public int insertAutoCreateTable(List<SubTableValue> subTableValues) {
        String sql = SqlSpeller.insertMultiTableMultiValuesUsingSuperTable(subTableValues);
        log.debug("SQL >>> " + sql);
        int affectRows = 0;
        try {
            affectRows = executeDDL(sql);
        } catch (SQLException e) {
            log.error("throws error: code:{},msg:{}", e.getErrorCode(), e.getMessage());
            e.printStackTrace();
        }
        return affectRows;
    }

    @Override
    public void updateSubTableTag(String database, String table, String tagName, String tagValue) {
        String tableName = StringUtils.hasText(database) ? database + "." + table : table;
        String sql = MessageFormat.format("ALTER TABLE {0} SET TAG {1}={2};", tableName, tagName, tagValue);
        log.debug("SQL >>> {}", sql);
        try {
            executeDDL(sql);
        } catch (SQLException e) {
            log.error("throws error: code:{},msg:{}", e.getErrorCode(), e.getMessage());
            e.printStackTrace();
        }
    }

}

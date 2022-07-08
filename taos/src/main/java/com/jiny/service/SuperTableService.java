package com.jiny.service;

import com.jiny.api.SuperTableInterface;
import com.jiny.entity.FieldMeta;
import com.jiny.entity.SuperTableMeta;
import com.jiny.entity.TagMeta;
import com.jiny.utils.SqlSpeller;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import org.springframework.util.StringUtils;

import java.sql.*;
import java.text.MessageFormat;
import java.util.ArrayList;
import java.util.List;

/**
 * @Auther: jiny
 * @CreateDate: 2022/7/8 16:56
 * @Description:
 */
@Service
@Slf4j
public class SuperTableService extends AbstractService implements SuperTableInterface {


    @Override
    public void create(SuperTableMeta superTableMeta) {
        String sql = SqlSpeller.createSuperTable(superTableMeta);
        log.debug("SQL >>> {}", sql);
        try {
            executeDDL(sql);
        } catch (SQLException e) {
            log.error("throws error: code:{},msg:{}", e.getErrorCode(), e.getMessage());
            e.printStackTrace();
        }
    }

    @Override
    public void drop(String database, String name) {
        String append = StringUtils.hasText(database) ? database + "." + name : name;
        String sql = "drop stable if exists " + append;
        log.debug("SQL >>> {}", sql);
        try {
            executeDDL(sql);
        } catch (SQLException e) {
            log.error("throws error: code:{},msg:{}", e.getErrorCode(), e.getMessage());
            e.printStackTrace();
        }
    }

    @Override
    public SuperTableMeta show(String database, String name) {
        String append = StringUtils.hasText(database) ? database + "." + name : name;
        String sql = "describe " + append;
        try (Connection connection = getConnection();
             Statement statement = connection.createStatement()){
            ResultSet resultSet = statement.executeQuery(sql);
            ResultSetMetaData meta = resultSet.getMetaData();
            SuperTableMeta superTableMeta = new SuperTableMeta();
            superTableMeta.setDatabase(database);
            superTableMeta.setName(name);
            List<FieldMeta> fieldMetas = new ArrayList<>();
            List<TagMeta> tagMetas = new ArrayList<>();
            while (resultSet.next()) {
                String fName = "", type = "", length = "", note = "";
                for (int i = 1; i <= meta.getColumnCount(); i++) {
                    String columnLabel = meta.getColumnLabel(i);
                    String value = resultSet.getString(i);
                    if (columnLabel.equals("Field")) {
                        fName = value;
                    } else if (columnLabel.equals("Type")) {
                        type = value;
                    } else if (columnLabel.equals("Length")) {
                        length = value;
                    } else if (columnLabel.equals("Note")) {
                        note = value;
                    }
                }
                type = type + "(" + length + ")";
                if (note.contains("TAG")) {
                    tagMetas.add(new TagMeta(fName, type));
                } else {
                    fieldMetas.add(new FieldMeta(fName, type));
                }
            }
            superTableMeta.setFields(fieldMetas);
            superTableMeta.setTags(tagMetas);
            return superTableMeta;
        } catch (SQLException e) {
            log.error("throws error: code:{},msg:{}", e.getErrorCode(), e.getMessage());
            e.printStackTrace();
            return null;
        }
    }

    @Override
    public void addField(String database, String table, FieldMeta fieldMeta) {
        String tableName = StringUtils.hasText(database) ? database + "." + table : table;
        String sql = MessageFormat.format("ALTER STABLE {0} ADD COLUMN {1} {2};", tableName, fieldMeta.getName(), fieldMeta.getType());
        log.debug("SQL >>> {}", sql);
        try {
            executeDDL(sql);
        } catch (SQLException e) {
            log.error("throws error: code:{},msg:{}", e.getErrorCode(), e.getMessage());
            e.printStackTrace();
        }
    }

    @Override
    public void delField(String database,String table, String name) {
        String tableName = StringUtils.hasText(database) ? database + "." + table : table;
        String sql = MessageFormat.format("ALTER STABLE {0} DROP COLUMN {1};", tableName, name);
        log.debug("SQL >>> {}", sql);
        try {
            executeDDL(sql);
        } catch (SQLException e) {
            log.error("throws error: code:{},msg:{}", e.getErrorCode(), e.getMessage());
            e.printStackTrace();
        }
    }

    @Override
    public void addTag(String database, String table, TagMeta tagMeta) {
        String tableName = StringUtils.hasText(database) ? database + "." + table : table;
        String sql = MessageFormat.format("ALTER STABLE {0} ADD TAG {1} {2};", tableName, tagMeta.getName(), tagMeta.getType());
        log.debug("SQL >>> {}", sql);
        try {
            executeDDL(sql);
        } catch (SQLException e) {
            log.error("throws error: code:{},msg:{}", e.getErrorCode(), e.getMessage());
            e.printStackTrace();
        }
    }

    @Override
    public void delTag(String database, String table, String name) {
        String tableName = StringUtils.hasText(database) ? database + "." + table : table;
        String sql = MessageFormat.format("ALTER STABLE {0} DROP TAG {1};", tableName, name);
        log.debug("SQL >>> {}", sql);
        try {
            executeDDL(sql);
        } catch (SQLException e) {
            log.error("throws error: code:{},msg:{}", e.getErrorCode(), e.getMessage());
            e.printStackTrace();
        }
    }

    @Override
    public void updateTag(String database, String table, String oldTag, String newTag) {
        String tableName = StringUtils.hasText(database) ? database + "." + table : table;
        String sql = MessageFormat.format("ALTER STABLE {0} CHANGE TAG {1} {2};", tableName, oldTag, newTag);
        log.debug("SQL >>> {}", sql);
        try {
            executeDDL(sql);
        } catch (SQLException e) {
            log.error("throws error: code:{},msg:{}", e.getErrorCode(), e.getMessage());
            e.printStackTrace();
        }
    }
}

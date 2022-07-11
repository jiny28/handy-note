package com.jiny.api;

import com.jiny.entity.SubTableMeta;
import com.jiny.entity.SubTableValue;

import java.util.List;

/**
 * @Auther: jiny
 * @CreateDate: 2022/7/11 11:31
 * @Description:
 */
public interface SubTableInterface {

    void create(SubTableMeta subTableMeta);

    int insert(List<SubTableValue> subTableValues);

    int insertAutoCreateTable(List<SubTableValue> subTableValues);

    void updateSubTableTag(String database, String table, String tagName, String tagValue);


}

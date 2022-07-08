package com.jiny.api;

import com.jiny.entity.FieldMeta;
import com.jiny.entity.SuperTableMeta;
import com.jiny.entity.TagMeta;

/**
 * @Auther: jiny
 * @CreateDate: 2022/7/8 16:45
 * @Description:
 */
public interface SuperTableInterface {

    void create(SuperTableMeta superTableMeta);

    void drop(String database, String name);

    SuperTableMeta show(String database, String name);

    void addField(String database, String table, FieldMeta fieldMeta);

    void delField(String database,String table,String name);

    void addTag(String database,String table ,TagMeta tagMeta);

    void delTag(String database,String table,String name);

    void updateTag(String database,String sTable, String oldTag, String newTag);





}

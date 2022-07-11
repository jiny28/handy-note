package com.jiny;

import com.jiny.api.QueryInterface;
import com.jiny.api.SubTableInterface;
import com.jiny.entity.*;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;

import java.util.*;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.core.Is.is;

@SpringBootTest
class SubTableTest {

    @Autowired
    private SubTableInterface subTableInterface;


    @Test
    void create() {
        subTableInterface.create(new SubTableMeta("", "meters", "d001", new ArrayList<TagValue>() {{
            add(new TagValue<>("location", "123"));
            add(new TagValue<>("groupId", 123));
        }}, null));
    }

    // 存在sql语句过长问题
    @Test
    void insert() {
        List<SubTableValue> subTableValues = new ArrayList<>();
        Random random = new Random();
        for (int i = 0; i < 1; i++) {
            SubTableValue subTableValue = new SubTableValue();
            subTableValue.setDatabase("");
            subTableValue.setName("d001");
            List<RowValue> rowValues = new ArrayList<>();
            for (int j = 0; j < 10; j++) {
                List<FieldValue> fieldValues = new ArrayList<FieldValue>() {{
                    add(new FieldValue("ts", new Date().getTime()));
                    add(new FieldValue("current", random.nextFloat()));
                    add(new FieldValue("voltage", random.nextInt()));
                    add(new FieldValue("phase", random.nextFloat()));
                }};
                RowValue rowValue = new RowValue(fieldValues);
                rowValues.add(rowValue);
                try {
                    Thread.sleep(100);
                } catch (InterruptedException e) {
                    e.printStackTrace();
                }
            }
            subTableValue.setValues(rowValues);
            subTableValues.add(subTableValue);
        }
        int insert = subTableInterface.insert(subTableValues);
        assertThat(insert, is(subTableValues.size() * subTableValues.get(0).getValues().size()));
    }

    @Test
    void insertAutoCreateTable() {
        List<SubTableValue> subTableValues = new ArrayList<>();
        Random random = new Random();
        for (int i = 0; i < 1; i++) {
            SubTableValue subTableValue = new SubTableValue();
            subTableValue.setDatabase("");
            subTableValue.setName("d001");
            subTableValue.setSuperTable("meters");
            subTableValue.setTags(new ArrayList<TagValue>() {{
                add(new TagValue<>("location", "123"));
                add(new TagValue<>("groupId", 123));
            }});
            List<RowValue> rowValues = new ArrayList<>();
            for (int j = 0; j < 10; j++) {
                List<FieldValue> fieldValues = new ArrayList<FieldValue>() {{
                    add(new FieldValue("ts", new Date().getTime()));
                    add(new FieldValue("current", random.nextFloat()));
                    add(new FieldValue("voltage", random.nextInt()));
                    add(new FieldValue("phase", random.nextFloat()));
                }};
                RowValue rowValue = new RowValue(fieldValues);
                rowValues.add(rowValue);
                try {
                    Thread.sleep(100);
                } catch (InterruptedException e) {
                    e.printStackTrace();
                }
            }
            subTableValue.setValues(rowValues);
            subTableValues.add(subTableValue);
        }
        int insert = subTableInterface.insertAutoCreateTable(subTableValues);
        assertThat(insert, is(subTableValues.size() * subTableValues.get(0).getValues().size()));

    }


    @Test
    void updateSubTableTag() {
        // 客户端需要配置FQDN，host设置映射到容器id或者是容器hostname
        subTableInterface.updateSubTableTag("", "d001", "groupid", "444");
    }

    @Autowired
    private QueryInterface queryInterface;


    @Test
    void select() {
        List<LinkedHashMap<String, String>> list = queryInterface.executeQuery("select * from meters");
        list.stream().forEach(map -> System.out.println(map));

    }

}

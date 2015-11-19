# user
+--------------------------+------------------+------+-----+---------+----------------+
| Field                    | Type             | Null | Key | Default | Extra          |
+--------------------------+------------------+------+-----+---------+----------------+
| user_id                  | int(10) unsigned | NO   | PRI | NULL    | auto_increment |
| user_name                | varbinary(255)   | NO   | UNI |         |                |
| user_real_name           | varbinary(255)   | NO   |     |         |                |
| user_password            | tinyblob         | NO   |     | NULL    |                |
| user_newpassword         | tinyblob         | NO   |     | NULL    |                |
| user_newpass_time        | binary(14)       | YES  |     | NULL    |                |
| user_email               | tinyblob         | NO   | MUL | NULL    |                |
| user_touched             | binary(14)       | NO   |     |         |                |
| user_token               | binary(32)       | NO   |     |         |                |
| user_email_authenticated | binary(14)       | YES  |     | NULL    |                |
| user_email_token         | binary(32)       | YES  | MUL | NULL    |                |
| user_email_token_expires | binary(14)       | YES  |     | NULL    |                |
| user_registration        | binary(14)       | YES  |     | NULL    |                |
| user_editcount           | int(11)          | YES  |     | NULL    |                |
| user_password_expires    | varbinary(14)    | YES  |     | NULL    |                |
+--------------------------+------------------+------+-----+---------+----------------+

# page https://www.mediawiki.org/wiki/Manual:Page_table
+--------------------+---------------------+------+-----+----------------+----------------+
| Field              | Type                | Null | Key | Default        | Extra          |
+--------------------+---------------------+------+-----+----------------+----------------+
| page_id            | int(10) unsigned    | NO   | PRI | NULL           | auto_increment |
| page_namespace     | int(11)             | NO   | MUL | NULL           |                |
| page_title         | varbinary(255)      | NO   |     | NULL           |                |
| page_is_redirect   | tinyint(3) unsigned | NO   | MUL | 0              |                |
| page_latest        | int(10) unsigned    | NO   |     | NULL           |                | points to  revision
| page_len           | int(10) unsigned    | NO   | MUL | NULL           |                |
+--------------------+---------------------+------+-----+----------------+----------------+

## todo
  - Add relationship to revision
  - Add nice page_title
  - add no duplicates for page_title
  - add primary key for page_title

# revision https://www.mediawiki.org/wiki/Manual:Revision_table
+--------------------+---------------------+------+-----+----------------+----------------+
| Field              | Type                | Null | Key | Default        | Extra          |
+--------------------+---------------------+------+-----+----------------+----------------+
| rev_id             | int(10) unsigned    | NO   | PRI | NULL           | auto_increment |
| rev_page           | int(10) unsigned    | NO   | MUL | NULL           |                | id points to page
| rev_text_id        | int(10) unsigned    | NO   |     | NULL           |                | id points to text
| rev_comment        | tinyblob            | NO   |     | NULL           |                |
| rev_user           | int(10) unsigned    | NO   | MUL | 0              |                |
| rev_user_text      | varbinary(255)      | NO   | MUL |                |                |
| rev_minor_edit     | tinyint(3) unsigned | NO   |     | 0              |                |
| rev_deleted        | tinyint(3) unsigned | NO   |     | 0              |                |
| rev_len            | int(10) unsigned    | YES  |     | NULL           |                |
| rev_parent_id      | int(10) unsigned    | YES  |     | NULL           |                |
| rev_sha1           | varbinary(32)       | NO   |     |                |                |
+--------------------+---------------------+------+-----+----------------+----------------+

When inserted update page_latest with link to this id.

# text https://www.mediawiki.org/wiki/Manual:Text_table
+-----------+------------------+------+-----+---------+----------------+
| Field     | Type             | Null | Key | Default | Extra          |
+-----------+------------------+------+-----+---------+----------------+
| old_id    | int(10) unsigned | NO   | PRI | NULL    | auto_increment |
| old_text  | mediumblob       | NO   |     | NULL    |                |
+-----------+------------------+------+-----+---------+----------------+

Need to parse the text from markdown. Then put it into this table. Insert a revision that points
to this text table. Then update the page_latest with the revision number.


# Steps
1. insert text
2. Try to get page or create new page
3. Insert revision
4. Update page with revision id
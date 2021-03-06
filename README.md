
# Git2Sqlite

For this project I wanted to hit a few birds with one stone:

1. Write something in Go, to learn the language.
2. Get an even deeper understanding of how Git internally works.
3. Produce something that might actually be useful.

Enter git2sqlite.

Git2Sqlite takes a target Git repository, reads up the object
database, and turns it into an sqlite database for future 
interrogation.

## Usage

After installing this Go command, you can move into a git repository 
(the top-level directory with the .git in it) and run just

``` git2sqlite ```

To execute the command. You can also provide a path:

``` git2sqlite path/to/repo ```

The command will output a file with the name <repository-name>.db into
the current working directory.

## SQLite Schema

```sql
CREATE TABLE commits (
      commit_hash text not null primary key,
      parent_a text,
      parent_b text,
      author text not null,
      committer text not null,
      date number not null,
      tree_hash not null,
      message text
    );
CREATE TABLE trees (
      tree_hash text not null, 
      blob_hash text not null,
      file_name text not null,
      file_mode text not null
    );
CREATE TABLE blobs (
      blob_hash text not null primary key, 
      content text
    );
CREATE TABLE heads (
    name text not null primary key, hash text
    );
CREATE TABLE tags (
    name text not null primary key, hash text
    );
```

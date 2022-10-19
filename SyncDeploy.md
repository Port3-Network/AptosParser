# Syncer
## 1. Clone code
```
git clone https://github.com/Port3-Network/AptosParser.git -b main # via http
git clone git@github.com:Port3-Network/AptosParser.git -b main # via ssh
```

## 2. Installing go dependencies
```
cd AptosParser && go mod tidy
```

## 2.1 go private library error resolution
<details>
<summary> fatal: could not read Username for 'https://github.com': terminal prompts disabled</summary>
github password generation: <a href="https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token">click here</a>

```bash
echo "machine github.com login <github email> password <github password>" > ~/.netrc
```
</details>

## 3. Edit configuration file
```
cp etc/config.conf.example etc/config.conf
```

## 3. build
```
cd aptos_sync && go build
```

## 4. run
```
./aptos_sync -n main
```


### 5. fix data to loog
```bash
select @@global.sql_mode;
+--------------------------------------------+
| @@global.sql_mode                          |
+--------------------------------------------+
| STRICT_TRANS_TABLES,NO_ENGINE_SUBSTITUTION |
+--------------------------------------------+
```
```bash

```
### 基本使用

```
func main() {
	var Conf Config

	configPath := flag.String("f", "./BurntSushi/toml/config/config.toml", "config file")
	flag.Parse()

	tomlDataByte, err := ioutil.ReadFile(*configPath)
	if err != nil {
		fmt.Println(err.Error())
	}

	if _, err := toml.Decode(string(tomlDataByte), &Conf); err != nil {
		fmt.Println(err.Error())
	}
	log.Println(Conf)

}
```

### [包原始地址](https://github.com/BurntSushi/toml)
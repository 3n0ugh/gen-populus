# gen-populus
<div align="center">
  <img src="./images/GSlogo.jpg" alt="drawing" width="180" height="160"/>
</div>

<p align="center"> Gen-Populus is a random population data generator.</p>


## Data Format
```Bash
   ID        --> Snowflake id generator [github.com/3n0ugh/snowflake]
   Name      --> Randomly chosen from name CSV file 
   Lastname  --> Randomly chosen from lastname CSV file
   Email     --> Name + Lastname + Last four number of ID + @3n0ugh.com
   Age       --> Generate random between 0 and 111
   Birthdate --> Random day and month + (Current year - age)
   Gender    --> It depends on which file the name is taken from.
```

## Requirements
- [Go (1.18)](https://go.dev/dl/)

## Usage

- You can define the output file, first name file, last name file, and population size. 
- However, the ratios such as the number of children-young-old or male-female numbers are defined randomly. 
(Check the [config](./internal/config/config.go) file to find out how. )
```go
// Open or if not exists create output file
file, _ := os.OpenFile("data.csv", os.O_CREATE|os.O_WRONLY, 0644)
cfg, _ := config.NewConfig(
	1e7,                               // population size: 10_000_000
	"./internal/data/female_name.csv", // female_name's file
	"./internal/data/male_name.csv",   // male_name's file
	"./internal/data/lastname.csv",    // lastname's file
	file)                              // output file
```
- Then, pass config as a parameter to the generator function.
```go
err = generator.Generate(cfg)
if err != nil {
	log.Println(err)
}
```
- You can check the [main](main.go) file for example.

## Benchmark
- Run benchmark (with 10 million population size):
```Bash
$ go test -bench=. -count=10  ./pkg/generator   
```
build:
	@go build -o suha.exe main.go

run: build
	suha $(arg)

test: build
	suha https://wallpaperaccess.com/full/682452.jpg

test_name: build
	suha https://i.pinimg.com/564x/75/ca/84/75ca84ce0204fe7d51ff8d333d857adc.jpg

test_utf: build
	suha https://blogger.googleusercontent.com/img/a/AVvXsEh5irh01lJQSd3_aRcPlDJLMG8XsU6XknJzcyEzi58JjFhuyRy5AmGy1bCkZ-itZQ4p1KE6bA8ET8kwiTylTYk8O6V75qoHeeg4q7KiAIqTT1phw1EoeyMbW5u706prkUdrZSfs56O3JWzR3Nx1cNnbJInU9Q3xLhVjuXp9VcDlUef3tFHfTBYZGwYs=s16000-rw

test_watch: build
	suha -w
test_dir:build
	suha https://i.pinimg.com/564x/75/ca/84/75ca84ce0204fe7d51ff8d333d857adc.jpg -o output
test_dir_aps:build
	suha https://i.pinimg.com/564x/75/ca/84/75ca84ce0204fe7d51ff8d333d857adc.jpg -o D:/aps
test_dir_watch:build
	suha -w -o output -f
test_auto: build
	suha -w -o output_auto -f -a
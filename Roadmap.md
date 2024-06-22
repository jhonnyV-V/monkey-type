## to add
- [x] typeOf
- [x] pop
- [x] reverse
- [x] join
- [x] split
- [x] replace
- [x] toLower
- [x] toUpper
- [x] trim
- [x] trimLeft
- [x] trimRight
- [x] contains and if is a string compare the runes
- [x] merge
- [x] findIndex
- [ ] lastIndex
- [ ] set and remove for hashmap
- [ ] for loop
- [ ] while loop

## maybe
- [ ] read code from file and execute?
- [ ] http client?
- [ ] file system access?

## Note for me

do a traditional for loop, the "parameters" should be a let statement, a condition and some expression
then the block, at compile time use the jump opcodes that exist

this should be something like this

- 1. create a new symbol table
- 2. execute the let statement
- 3. save this point to jump later on
- 4. check the condition and use OpJumpNotTruthy
- 5. execute the code in the block
- 6. execute the expression
- 7. jump to point 3


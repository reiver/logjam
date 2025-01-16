# LOGJAM lib/

The `lib/` contains libraries that **should never** `import` anything from `logjam` outside of the `lib/` directory.

The point is that this is decoupled from the rest of this code-base.
And the rest of this code-base depends on it, _not_ vice versa.

from archlinux as builder

run pacman -Syu --noconfirm --needed go gcc

copy . /build

workdir /build
run go build -o main -trimpath .

from archlinux

run pacman -Syu --noconfirm

copy --from=builder /build/main /app/server

entrypoint ["/app/server"]

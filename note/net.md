# package net 的笔记

## 参考

https://golang.org/pkg/net/

## 常量

```go
const (
  IPv4len = 4
  IPv6len = 16
)
```

## 变量

```go
var (
  IPv4bcast     = IPv4(255, 255, 255, 255) // limited broadcast
  IPv4allsys    = IPv4(224, 0, 0, 1)       // all systems
  IPv4allrouter = IPv4(224, 0, 0, 2)       // all routers
  IPv4zero      = IPv4(0, 0, 0, 0)         // all zeros
)

var (
  IPv6zero                   = IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
  IPv6unspecified            = IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
  IPv6loopback               = IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
  IPv6interfacelocalallnodes = IP{0xff, 0x01, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x01}
  IPv6linklocalallnodes      = IP{0xff, 0x02, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x01}
  IPv6linklocalallrouters    = IP{0xff, 0x02, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x02}
)
```

## 结构体

Addr 表示一个网络端点地址.

```go
type Addr interface {
  Network() string // name of the network (for example, "tcp", "udp")
  String() string  // string form of address (for example, "192.0.2.1:25", "[2001:db8::1]:80")
}
```

Buffers 包含零个或多个要写入的 bytes.

```go
type Buffers [][]byte
func (v *Buffers) Read(p []byte) (n int, err error)
func (v *Buffers) WriteTo(w io.Writer) (n int64, err error)
```

---

**Conn** 是一个通用的面向数据流(stream-oriented)的网络连接.
多个 goroutine 可以同时在一个 Conn 上调用方法.

```go
type Conn interface {
  // Read 从连接中读取数据.
  // Read 可以使用超时机制并在固定的时间到达时返回一个 Error 当 Timeout() == true
  // see SetDeadline and SetReadDeadline.
  Read(b []byte) (n int, err error)

  // Write 将数据写入到连接中.
  // Write 可以使用超时机制并在固定的时间到达时返回一个 Error 当 Timeout() == true
  // see SetDeadline and SetWriteDeadline.
  Write(b []byte) (n int, err error)

  // Close 关闭连接.
  // 任何 blocked 读或写操作将被 unblocked 并返回错误.
  Close() error

  // LocalAddr 返回本地网络地址.
  LocalAddr() Addr

  // RemoteAddr 返回远程网络地址.
  RemoteAddr() Addr

  // SetDeadline 在连接上设置读写截止时间.
  // 它等价于同时调用
  // SetReadDeadline and SetWriteDeadline.
  //
  // 一个截止时间是一个绝对时间,  在那之后 I/O 操作
  // 会返回 timeout Error 并失败, 而不是继续阻塞.
  // 截止时间会用在所有未来或待定状态的 I/O 上,
  // 而不是紧随其后的 Read 或 Write 上.
  // 当截止时间被执行时, 连接可以通过设置一个新的未来的截止时间来拒绝.
  //
  // 一个 idle 超时可以被实现为, 在成功读取或写入之后不断延长截止时间.
  //
  // t 的零值表示 I/O 操作永远不会超时
  //
  // 注意 如果 TCP 连接打开了 keep-alive,
  // 这是默认的行为, 除非通过 Dialer.KeepAlive 或
  // ListenConfig.KeepAlive 覆盖,
  // 然后一个 keep-alive 失败也可能返回一个超时错误.
  // 在 Unix 系统上, 一个 keep-alive 失败可以被
  // errors.Is(err, syscall.ETIMEDOUT) 检测到.
  SetDeadline(t time.Time) error

  // SetReadDeadline 在未来的 Read 和任何当前阻塞的 Read 上设置超时.
  // t 的零值表示 Read 操作永远不会超时
  SetReadDeadline(t time.Time) error

  // SetWriteDeadline 在未来的 Write 和任何当前阻塞的 Write 上设置超时.
  // 即使写入超时了, 它也可能返回 n > 0,
  // 表示一些数据已经成功写入了.
  // t 的零值表示 Write 操作永远不会超时
  SetWriteDeadline(t time.Time) error
}
```

**Dial** 在指定的 network 上连接 address.

```go
func Dial(network, address string) (Conn, error)
```

Known networks are "tcp", "tcp4" (IPv4-only), "tcp6" (IPv6-only), "udp", "udp4" (IPv4-only), "udp6" (IPv6-only), "ip", "ip4" (IPv4-only), "ip6" (IPv6-only), "unix", "unixgram" and "unixpacket".

对于 TCP 和 UDP 网络, address 的结构是 `host:port`. host 必须是一个字面量形式的 IP 地址, 或者一个可以被
解析为 IP 地址的主机名. port 必须是看一个字面量形式的端口或者一个服务名. 如果 host 是一个 IPv6 的字面量, 那么
必须用 `[]` 引起来, 比如 `[2001:db8::1]:80` 和 `[fe80::1%zone]:80`. zone 指示了字面量形式的 IPv6 地址的
范围, 定义在 `RFC 4007`. JoinHostPort 和 SplitHostPort 函数操纵这种形式的 host 和 port 对. 当使用 TCP,
且 host 可以被解析为多个 IP 地址, Dial 会按顺序尝试它们, 直到其中一个成功了.

```go
Dial("tcp", "golang.org:http")
Dial("tcp", "192.0.2.1:http")
Dial("tcp", "198.51.100.1:80")
Dial("udp", "[2001:db8::1]:domain")
Dial("udp", "[fe80::1%lo0]:53")
Dial("tcp", ":80")
```

对于 IP 网络, network 必须是 `ip`, `ip4`, `ip6` 接着一个冒号, 一个协议数字或者一个协议名字,
address 的结构和 host 一样. host 必须是一个字面量形式的 IP 地址, 或者一个带有 zone 的
字面量形式的 IPv6 地址. 这取决于操作系统在非知名的协议数字上的表现, 比如 "0" 或 "255".

```go
Dial("ip4:1", "192.0.2.1")
Dial("ip6:ipv6-icmp", "2001:db8::1")
Dial("ip6:58", "fe80::1%lo0")
```

对于 TCP, UDP 和 IP 网络, 如果 host 是空的或者字面量形式没有指定 IP 地址, 比如 TCP 和 UDP 的
`:80`, `0.0.0.0:80` 或 `[::]:80`, IP 的 `0.0.0.0` 和 `::`, 假定为本地系统.

对于 Unix 网络, address 必须是一个文件系统路径.

```go
func DialTimeout(network, address string, timeout time.Duration) (Conn, error)
func FileConn(f *os.File) (c Conn, err error)
```

FileConn 返回一个在打开的文件 f 上的网络连接的副本. 调用者有责任在最后关闭 f.
关闭 c 不会影响 f, 关闭 f 也不会影响 c.

---

**Dialer** 包含连接到一个地址的选项.

每一个字段都是零值等价于不使用选项. 使用零值的 Dialer 相当于直接调用 Dial 函数.

```go
type Dialer struct {
  // Timeout 是一个 dial 等待连接完成的最大时间.
  // 如果 Deadline 也设置了, 可能会提前失败.
  //
  // 默认是没有超时.
  //
  // 当使用 TCP 且 dialing 一个主机名有多个 IP 地址,
  // 超时会除以它们的数量.
  //
  // 不管是否使用超时, 操作系统可能强迫他更早结束.
  // 比如 TCP 超时通常在 3 跟中的时候.
  Timeout time.Duration

  // Deadline 是一个 dials 会失败的绝对时间.
  // 如果 Timeout 也设置了, 可能更早失败.
  // 零值表示没有截止时间, 或者和 Timeout 一样取决于操作系统.
  Deadline time.Time

  // LocalAddr 是本地地址, 用于 dial 一个 address 时.
  // 这个地址的类型必须兼容要 dial 的地址.
  // 如果是 nil, 会自动选择一个本地地址.
  LocalAddr Addr

  // DualStack 预先启用 RFC 6555 Fast Fallback
  // 支持, 又称 "Happy Eyeballs".
  // 当 IPv6 可能配置错误或者挂起时,
  // 会尽快尝试 IPv4.
  //
  // Deprecated: Fast Fallback 是默认启用的.
  // 要禁止它, 将 FallbackDelay 设置为一个负值.
  DualStack bool // Go 1.2

  // FallbackDelay 指定了等待的时长, 直到开始
  // spawning a RFC 6555 Fast Fallback connection.
  // 换句话说, 这是等待 IPv6 成功连接的时间,
  // 超时后会假设 IPv6 配置错误, 并回退到 IPv4.
  //
  // 如果是 0, 默认的延迟 300ms 会被使用.
  // 一个负值会禁用 Fast Fallback 支持.
  FallbackDelay time.Duration // Go 1.5

  // KeepAlive 指定了一个间隔事件, 用于 keep-alive
  // 探针检测活动的网络连接.
  // 如果是 0, keep-alive 探针被设置为默认值,
  // (当前是 15 秒), 如果协议和操作系统支持的话.
  // 不支持的协议和操作系统会忽略这个字段.
  // 如果是负值, 禁用 keep-alive probes
  KeepAlive time.Duration // Go 1.3

  // Resolver 可选地用于指定替代的 resolver.
  Resolver *Resolver // Go 1.8

  // Cancel 是一个可选的 channel, 它的关闭会导致
  // dial 应该被取消. 不是所有的 dial 都支持取消.
  //
  // Deprecated: 使用 DialContext 代替.
  Cancel <-chan struct{} // Go 1.6

  // 如果 Control 不是 nil, 它被在创建网络连接但还未实际 dial 时被调用.
  //
  // 网络和地址操作会被传递给 Control 方法, 但不一定是传递给 Dial 的那些.
  // 比如, 传递 "tcp" 给 Dial 会导致 Control 会通过 "tcp4" 或 "tcp6" 调用.
  Control func(network, address string, c syscall.RawConn) error // Go 1.11
}
```

Dialer 相关的方法.

```go
func (d *Dialer) Dial(network, address string) (Conn, error)
func (d *Dialer) DialContext(ctx context.Context, network, address string) (Conn, error)
```

DialContext 使用提供的 context 连接到指定的网络的地址上.

提供的 Context 不能为 nil. 如果 context 在连接完成前就过期了, 会返回一个错误.
一旦成功连接, context 的任何过期都不会影响连接.

当使用 TCP, 且 address 参数中的 host 被解析为多个网络地址, 任何 dial 超时(来自 d.Timeout 或 ctx),
会被分割到每个连续的 dial 上, 每一个都会被给予一个恰当的时间来连接. 比如, 一个 host 有 4 个 IP 地址,
超时是一分钟, 每个地址都会有 15 秒的时间来尝试完成连接.

---

IPConn 实现了 Conn 和 PacketConn 接口, 用于 IP 网络连接.

```go
type IPConn struct {
  // contains filtered or unexported fields
}
```

相关的方法和函数.

```go
// 连接到服务器
func DialIP(network string, laddr, raddr *IPAddr) (*IPConn, error)
// 本地监听
func ListenIP(network string, laddr *IPAddr) (*IPConn, error)
func (c *IPConn) Close() error
func (c *IPConn) File() (f *os.File, err error)
func (c *IPConn) LocalAddr() Addr
func (c *IPConn) Read(b []byte) (int, error)
func (c *IPConn) ReadFrom(b []byte) (int, Addr, error)
func (c *IPConn) ReadFromIP(b []byte) (int, *IPAddr, error)
func (c *IPConn) ReadMsgIP(b, oob []byte) (n, oobn, flags int, addr *IPAddr, err error)
func (c *IPConn) RemoteAddr() Addr
func (c *IPConn) SetDeadline(t time.Time) error
func (c *IPConn) SetReadBuffer(bytes int) error
func (c *IPConn) SetReadDeadline(t time.Time) error
func (c *IPConn) SetWriteBuffer(bytes int) error
func (c *IPConn) SetWriteDeadline(t time.Time) error
func (c *IPConn) SyscallConn() (syscall.RawConn, error)
func (c *IPConn) Write(b []byte) (int, error)
func (c *IPConn) WriteMsgIP(b, oob []byte, addr *IPAddr) (n, oobn int, err error)
func (c *IPConn) WriteTo(b []byte, addr Addr) (int, error)
func (c *IPConn) WriteToIP(b []byte, addr *IPAddr) (int, error)
```

---

IPNet 表示一个 IP 网络.

```go
type IPNet struct {
  IP   IP     // network number
  Mask IPMask // network mask
}
```

Interface 表示网络接口名字和索引的映射, 也表示网络接口设备的信息.

```go
type Interface struct {
  Index        int          // positive integer that starts at one, zero is never used
  MTU          int          // maximum transmission unit
  Name         string       // e.g., "en0", "lo0", "eth0.100"
  HardwareAddr HardwareAddr // IEEE MAC-48, EUI-48 and EUI-64 form
  Flags        Flags        // e.g., FlagUp, FlagLoopback, FlagMulticast
}
```

---

ListenConfig 包含监听一个地址时的选项.

```go
type ListenConfig struct {
  // 如果 Control 不是 nil, 它会在创建网络连接之后, 但在绑定到操作系统之前被调用.
  //
  // Network 和 address 参数被传递给 Control 方法的, 不一定是传递给 Listen 的那些.
  // 举个例子, 传递 "tcp" 给 Listen 会导致 Control 函数接收到 "tcp4" 或 "tcp6".
  Control func(network, address string, c syscall.RawConn) error

  // KeepAlive 指定这个 listener 用于网络连接的 keep-alive 时间.
  // 如果是零, keep-alives 会被启动, 如果协议和操作系统也支持的话.
  // 网络协议或操作系统不支持的话, 会忽略这个字段.
  // 如果是负数, 禁用 keep-alives
  KeepAlive time.Duration // Go 1.13
}
```

```go
func (lc *ListenConfig) Listen(ctx context.Context, network, address string) (Listener, error)
func (lc *ListenConfig) ListenPacket(ctx context.Context, network, address string) (PacketConn, error)
```

---

**Listener** 是一个通用的面向流式协议(stream-oriented protocols)的网络监听器.

多个 goroutine 可能在一个监听器上同时调用方法.

```go
type Listener interface {
  // Accept 等待并返回下一个 listener 的连接.
  Accept() (Conn, error)

  // Close 关闭 listener.
  // 任何阻塞的 Accept 操作会被释放且返回错误.
  Close() error

  // Addr 返回监听器的网络地址.
  Addr() Addr
}
```

```go
func FileListener(f *os.File) (ln Listener, err error)
func Listen(network, address string) (Listener, error)
```

Listen 监听一个本地网络地址.

network 参数必须是 "tcp", "tcp4", "tcp6", "unix" 或 "unixpacket".

---

PacketConn 是一个通用的面向包(packet-oriented)的网络连接.

多个 goroutine 可能在一个 PacketConn 上同时调用方法.

```go
type PacketConn interface {
  // ReadFrom reads a packet from the connection,
  // copying the payload into p. It returns the number of
  // bytes copied into p and the return address that
  // was on the packet.
  // It returns the number of bytes read (0 <= n <= len(p))
  // and any error encountered. Callers should always process
  // the n > 0 bytes returned before considering the error err.
  // ReadFrom can be made to time out and return
  // an Error with Timeout() == true after a fixed time limit;
  // see SetDeadline and SetReadDeadline.
  ReadFrom(p []byte) (n int, addr Addr, err error)

  // WriteTo writes a packet with payload p to addr.
  // WriteTo can be made to time out and return
  // an Error with Timeout() == true after a fixed time limit;
  // see SetDeadline and SetWriteDeadline.
  // On packet-oriented connections, write timeouts are rare.
  WriteTo(p []byte, addr Addr) (n int, err error)

  // Close closes the connection.
  // Any blocked ReadFrom or WriteTo operations will be unblocked and return errors.
  Close() error

  // LocalAddr returns the local network address.
  LocalAddr() Addr

  // SetDeadline sets the read and write deadlines associated
  // with the connection. It is equivalent to calling both
  // SetReadDeadline and SetWriteDeadline.
  //
  // A deadline is an absolute time after which I/O operations
  // fail with a timeout (see type Error) instead of
  // blocking. The deadline applies to all future and pending
  // I/O, not just the immediately following call to ReadFrom or
  // WriteTo. After a deadline has been exceeded, the connection
  // can be refreshed by setting a deadline in the future.
  //
  // An idle timeout can be implemented by repeatedly extending
  // the deadline after successful ReadFrom or WriteTo calls.
  //
  // A zero value for t means I/O operations will not time out.
  SetDeadline(t time.Time) error

  // SetReadDeadline sets the deadline for future ReadFrom calls
  // and any currently-blocked ReadFrom call.
  // A zero value for t means ReadFrom will not time out.
  SetReadDeadline(t time.Time) error

  // SetWriteDeadline sets the deadline for future WriteTo calls
  // and any currently-blocked WriteTo call.
  // Even if write times out, it may return n > 0, indicating that
  // some of the data was successfully written.
  // A zero value for t means WriteTo will not time out.
  SetWriteDeadline(t time.Time) error
}
```

```go
func FilePacketConn(f *os.File) (c PacketConn, err error)
func ListenPacket(network, address string) (PacketConn, error)
```

network 参数必须是 "udp", "udp4", "udp6", "unixgram", 或一个 IP 负载.

---

Resolver 寻找名字和数字. 所有的 nil \*Resolver 等价于 Resolver 的零值.

```go
type Resolver struct {
  // PreferGo 控制 Go 内置的 DNS resolver 是否是首选的, 如果它在该系统上可用.
  // 等价于设置 GODEBUG=netdns=go, 但仅限于这个解析器.
  PreferGo bool

  // StrictErrors 控制临时错误的行为
  // (包括超时, socket 错误, 和 SERVFAIL), 当使用 Go 内置的解析器时.
  // 对于一个由多个子查询组成的查询 (比如 A+AAAA 地址查询, 或者遍历 DNS 搜索列表),
  // 这个选项制造错误来禁止整个查询, 而不是返回部分结果.
  // 默认它是不启用的, 因为它可能影响兼容性, 解析器不能正确处理 AAAA 查询.
  StrictErrors bool // Go 1.9

  // Dial 可选地指定一个替代的 dialer 用于 Go 内置的 DNS 解析器的 TCP 和 UDP 连接.
  // host 和 address 参数总是一个字面量形式的 IP 地址, 而不是一个主机名,
  // address 参数的 port 总是一个字面量形式的端口号, 而不是一个服务名字.
  // 如果 Conn 返回一个 PacketConn, 发送和接收 DNS 信息必须符合
  // RFC 1035 section 4.2.1, "UDP usage".
  // 否者, 负载于 Conn 上的 DNS 信息必须符合
  // RFC 7766 section 5, "Transport Protocol Selection".
  // 如果是 nil, 默认的 dialer 会被使用.
  Dial func(ctx context.Context, network, address string) (Conn, error) // Go 1.9
  // contains filtered or unexported fields
}
```

```go
func (r *Resolver) LookupAddr(ctx context.Context, addr string) (names []string, err error)
func (r *Resolver) LookupCNAME(ctx context.Context, host string) (cname string, err error)
func (r *Resolver) LookupHost(ctx context.Context, host string) (addrs []string, err error)
func (r *Resolver) LookupIPAddr(ctx context.Context, host string) ([]IPAddr, error)
func (r *Resolver) LookupMX(ctx context.Context, name string) ([]*MX, error)
func (r *Resolver) LookupNS(ctx context.Context, name string) ([]*NS, error)
func (r *Resolver) LookupPort(ctx context.Context, network, service string) (port int, err error)
func (r *Resolver) LookupSRV(ctx context.Context, service, proto, name string) (cname string, addrs []*SRV, err error)
func (r *Resolver) LookupTXT(ctx context.Context, name string) ([]string, error)
```

---

TCPConn 在 TCP 网络连接上实现了 Conn 接口.

```go
type TCPConn struct {
  //tered or unexported fields
}
```

```go
func DialTCP(network string, laddr, raddr *TCPAddr) (*TCPConn, error)
func (c *TCPConn) Close() error
func (c *TCPConn) CloseRead() error
func (c *TCPConn) CloseWrite() error
func (c *TCPConn) File() (f *os.File, err error)
func (c *TCPConn) LocalAddr() Addr
func (c *TCPConn) Read(b []byte) (int, error)
func (c *TCPConn) ReadFrom(r io.Reader) (int64, error)
func (c *TCPConn) RemoteAddr() Addr
func (c *TCPConn) SetDeadline(t time.Time) error
func (c *TCPConn) SetKeepAlive(keepalive bool) error
func (c *TCPConn) SetKeepAlivePeriod(d time.Duration) error
func (c *TCPConn) SetLinger(sec int) error
func (c *TCPConn) SetNoDelay(noDelay bool) error
func (c *TCPConn) SetReadBuffer(bytes int) error
func (c *TCPConn) SetReadDeadline(t time.Time) error
func (c *TCPConn) SetWriteBuffer(bytes int) error
func (c *TCPConn) SetWriteDeadline(t time.Time) error
func (c *TCPConn) SyscallConn() (syscall.RawConn, error)
func (c *TCPConn) Write(b []byte) (int, error)
```

```go
type TCPListener struct {
  // contains filtered or unexported fields
}
func ListenTCP(network string, laddr *TCPAddr) (*TCPListener, error)
func (l *TCPListener) Accept() (Conn, error)
func (l *TCPListener) AcceptTCP() (*TCPConn, error)
func (l *TCPListener) Addr() Addr
func (l *TCPListener) Close() error
func (l *TCPListener) File() (f *os.File, err error)
func (l *TCPListener) SetDeadline(t time.Time) error
func (l *TCPListener) SyscallConn() (syscall.RawConn, error)
```

---

UDPConn 在 UDP 网络连接上实现了 Conn 和 PacketConn 接口.

```go
type UDPConn struct {
  // contains filtered or unexported fields
}
```

```go
func DialUDP(network string, laddr, raddr *UDPAddr) (*UDPConn, error)
func ListenMulticastUDP(network string, ifi *Interface, gaddr *UDPAddr) (*UDPConn, error)
func ListenUDP(network string, laddr *UDPAddr) (*UDPConn, error)
func (c *UDPConn) Close() error
func (c *UDPConn) File() (f *os.File, err error)
func (c *UDPConn) LocalAddr() Addr
func (c *UDPConn) Read(b []byte) (int, error)
func (c *UDPConn) ReadFrom(b []byte) (int, Addr, error)
func (c *UDPConn) ReadFromUDP(b []byte) (int, *UDPAddr, error)
func (c *UDPConn) ReadMsgUDP(b, oob []byte) (n, oobn, flags int, addr *UDPAddr, err error)
func (c *UDPConn) RemoteAddr() Addr
func (c *UDPConn) SetDeadline(t time.Time) error
func (c *UDPConn) SetReadBuffer(bytes int) error
func (c *UDPConn) SetReadDeadline(t time.Time) error
func (c *UDPConn) SetWriteBuffer(bytes int) error
func (c *UDPConn) SetWriteDeadline(t time.Time) error
func (c *UDPConn) SyscallConn() (syscall.RawConn, error)
func (c *UDPConn) Write(b []byte) (int, error)
func (c *UDPConn) WriteMsgUDP(b, oob []byte, addr *UDPAddr) (n, oobn int, err error)
func (c *UDPConn) WriteTo(b []byte, addr Addr) (int, error)
func (c *UDPConn) WriteToUDP(b []byte, addr *UDPAddr) (int, error)
```

---

UnixConn 实现了 Conn 接口, 用于连接到 Unix domain sockets.

```go
type UnixConn struct {
  // contains filtered or unexported fields
}
```

```go
func DialUnix(network string, laddr, raddr *UnixAddr) (*UnixConn, error)
func ListenUnixgram(network string, laddr *UnixAddr) (*UnixConn, error)
func (c *UnixConn) Close() error
func (c *UnixConn) CloseRead() error
func (c *UnixConn) CloseWrite() error
func (c *UnixConn) File() (f *os.File, err error)
func (c *UnixConn) LocalAddr() Addr
func (c *UnixConn) Read(b []byte) (int, error)
func (c *UnixConn) ReadFrom(b []byte) (int, Addr, error)
func (c *UnixConn) ReadFromUnix(b []byte) (int, *UnixAddr, error)
func (c *UnixConn) ReadMsgUnix(b, oob []byte) (n, oobn, flags int, addr *UnixAddr, err error)
func (c *UnixConn) RemoteAddr() Addr
func (c *UnixConn) SetDeadline(t time.Time) error
func (c *UnixConn) SetReadBuffer(bytes int) error
func (c *UnixConn) SetReadDeadline(t time.Time) error
func (c *UnixConn) SetReadDeadline(t time.Time) error
func (c *UnixConn) SetWriteDeadline(t time.Time) error
func (c *UnixConn) SyscallConn() (syscall.RawConn, error)
func (c *UnixConn) Write(b []byte) (int, error)
func (c *UnixConn) WriteMsgUnix(b, oob []byte, addr *UnixAddr) (n, oobn int, err error)
func (c *UnixConn) WriteTo(b []byte, addr Addr) (int, error)
func (c *UnixConn) WriteToUnix(b []byte, addr *UnixAddr) (int, error)
```

---

UnixListener

```go
type UnixListener struct {
  // contains filtered or unexported fields
}
func ListenUnix(network string, laddr *UnixAddr) (*UnixListener, error)
func (l *UnixListener) Accept() (Conn, error)
func (l *UnixListener) AcceptUnix() (*UnixConn, error)
func (l *UnixListener) Addr() Addr
func (l *UnixListener) Close() error
func (l *UnixListener) File() (f *os.File, err error)
func (l *UnixListener) SetDeadline(t time.Time) error
func (l *UnixListener) SetUnlinkOnClose(unlink bool)
func (l *UnixListener) SyscallConn() (syscall.RawConn, error)
```

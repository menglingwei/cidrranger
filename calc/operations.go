package calc

import (
    "net"
    "encoding/binary"
    "math"
    //"strconv"
    "github.com/pkg/errors"
)

func (this *CalculatorService) CidrFromParsedIPs(ips []string) ([]*net.IPNet, error) {
    parsedIps := make([]net.IP, 0)

    for i, ip := range ips {
        pip := net.ParseIP(ip)

        if pip == nil {
            return nil, errors.Errorf("unsable to parse ips[%d] == %s", i, ip)
        }

        parsedIps = append(parsedIps, pip)
    }

    return this.CidrFromIPs(parsedIps), nil
}

func (this *CalculatorService) CidrFromIPs(ips []net.IP) []*net.IPNet {
    ns := make([]*net.IPNet, 0)
    upperLower := this.UpperLowerIP(ips)

    type xxxx struct {
        n     net.IPNet
        found bool
    }

    um := make(map[string]xxxx)
    lm := make(map[string]xxxx)

    for ones := 1; ones <= 32; ones++ {
        m := net.CIDRMask(ones, net.IPv4len * 8)
        unet := &net.IPNet{ IP: upperLower[1], Mask: m, }
        lnet := &net.IPNet{ IP: upperLower[0], Mask: m, }



        this.Logger.Printf("========> [%d] mask(hex): %s", ones, m)
        this.Logger.Printf("========> ipnet: %s", n.String())
        this.Logger.Printf("========> first: %s, last: %s", endpoints[0], endpoints[1])
        //this.Logger.Printf("========> mask: %s", strconv.FormatUint(imask, 2))

    }

    return ns
}

func (this *CalculatorService) NetworkAddressFromCidr(cidr string) (net.IP, error) {
    _, n, err := net.ParseCIDR(cidr)

    if err != nil {
        return nil, err
    }

    return this.NetworkAddress(n), nil
}

func (this *CalculatorService) NetworkAddress(n *net.IPNet) net.IP {
    return n.IP.Mask(n.Mask)
}

func (this *CalculatorService) BroadcastAddressFromCidr(cidr string) (net.IP, error) {
    _, n, err := net.ParseCIDR(cidr)

    if err != nil {
        return nil, err
    }

    return this.BroadcastAddress(n), nil
}

func (this *CalculatorService) BroadcastAddress(n *net.IPNet) net.IP {
    cidrMask := net.CIDRMask(n.Mask.Size())
    broadcast := net.IP(make([]byte, 4))

    //if this.Debug {
    //    this.Logger.Printf("len(broadcast): %d", len(broadcast))
    //    this.Logger.Printf("len(n.IP.To4()): %d", len(n.IP.To4()))
    //    this.Logger.Printf("len(n.IP): %d", len(n.IP))
    //    this.Logger.Printf("len(cidrMask): %d", len(cidrMask))
    //}

    for i, _ := range n.IP.To4() {
        broadcast[i] = n.IP[i] | ^cidrMask[i]
    }

    return broadcast
}

func (this *CalculatorService) IncIP(ip net.IP) net.IP {
    next := net.IP(make([]byte, len(ip)))
    copy(next, ip)

    for i := len(next) - 1; i >= 0; i-- {
        next[i]++
        if next[i] != 0 {
            break
        }
    }

    return next
}

func (this *CalculatorService) DecIP(ip net.IP) net.IP {
    prev := net.IP(make([]byte, len(ip)))
    copy(prev, ip)

    for i := len(prev) - 1; i >= 0; i-- {
        prev[i]--
        if prev[i] != 255 {
            break
        }
    }

    return prev
}

func (this *CalculatorService) SubnetEnumerateFromCidr(cidr string) ([]net.IP, error) {
    _, n, err := net.ParseCIDR(cidr)

    if err != nil {
        return nil, err
    }

    return this.SubnetEnumerate(n), nil
}

func (this *CalculatorService) SubnetEnumerate(n *net.IPNet) []net.IP {
    ips := make([]net.IP, 0)

    naddr := this.NetworkAddress(n)

    for ip := naddr; n.Contains(ip); ip = this.IncIP(ip) {
        ips = append(ips, ip)
    }

    if this.RemoveNetworkAddress {
        ips = ips[1:]
    }

    if this.RemoveBroadcastAddress {
        ips = ips[:len(ips) - 1]
    }

    return ips
}

func (this *CalculatorService) SubnetEndpointsFromCidr(cidr string) ([]net.IP, error) {
    _, n, err := net.ParseCIDR(cidr)

    if err != nil {
        return nil, err
    }

    return this.SubnetEndpoints(n), nil
}

func (this *CalculatorService) SubnetEndpoints(n *net.IPNet) []net.IP {
    ips := make([]net.IP, 2)

    naddr := this.NetworkAddress(n)

    if this.RemoveNetworkAddress {
        ips[0] = this.IncIP(naddr)
    } else {
        ips[0] = naddr
    }

    baddr := this.BroadcastAddress(n)

    if this.RemoveBroadcastAddress {
        ips[1] = this.DecIP(baddr)
    } else {
        ips[1] = baddr
    }

    return ips
}

func (this *CalculatorService) UpperLowerIP(ips []net.IP) []net.IP {
    h := uint32(0)
    l := uint32(math.MaxUint32)

    for _, ip := range ips {
        x := uint32(ip2int(ip))

        if x > h {
            h = x
        }

        if x < l {
            l = x
        }
    }

    return []net.IP{int2ip(l), int2ip(h)}
}

func (this *CalculatorService) UpperIP(ips []net.IP) net.IP {
    h := uint32(0)

    for _, ip := range ips {
        x := uint32(ip2int(ip))

        if x > h {
            h = x
        }
    }

    return int2ip(h)
}

func (this *CalculatorService) LowerIP(ips []net.IP) net.IP {
    l := uint32(math.MaxUint32)

    for _, ip := range ips {
        x := uint32(ip2int(ip))

        if x < l {
            l = x
        }
    }

    return int2ip(l)
}

//https://gist.github.com/ammario/649d4c0da650162efd404af23e25b86b
func ip2int(ip net.IP) uint32 {
    if len(ip) == 16 {
        return binary.BigEndian.Uint32(ip[12:16])
    }
    return binary.BigEndian.Uint32(ip)
}

//https://gist.github.com/ammario/649d4c0da650162efd404af23e25b86b
func int2ip(nn uint32) net.IP {
    ip := make(net.IP, 4)
    binary.BigEndian.PutUint32(ip, nn)
    return ip
}

const fs = require('fs')
const path = require('path')
const stream = require('stream')

console.log("# start")
if (!process.argv[2]) {
  return new Error("Lack args")
}
let e = process.argv[2]
let t
for (var a = {
    data: {
      info: e
    }
  }, s = {
    q: r,
    h: n,
    m: i,
    k: o
  }, l = a.data.info, u = l.substring(l.length - 4).split(""), c = 0; c < u.length; c++)
  u[c] = u[c].toString().charCodeAt(0) % 4;
u.reverse();
for (var d = [], c = 0; c < u.length; c++)
  d.push(l.substring(u[c] + 1, u[c] + 2)),
  l = l.substring(0, u[c] + 1) + l.substring(u[c] + 2);
a.data.encrypt_table = d,
  a.data.key_table = [];
for (var c in a.data.encrypt_table)
  "q" != a.data.encrypt_table[c] && "k" != a.data.encrypt_table[c] || (a.data.key_table.push(l.substring(l.length - 12)),
    l = l.substring(0, l.length - 12));
a.data.key_table.reverse(),
  a.data.info = l;
var h = new Array(-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 62, -1, -1, -1, 63, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, -1, -1, -1, -1, -1, -1, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, -1, -1, -1, -1, -1, -1, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, -1, -1, -1, -1, -1);
a.data.info = function (e) {
  var t, r, n, i, o, a, s;
  for (a = e.length,
    o = 0,
    s = ""; o < a;) {
    do {
      t = h[255 & e.charCodeAt(o++)]
    } while (o < a && -1 == t);
    if (-1 == t)
      break;
    do {
      r = h[255 & e.charCodeAt(o++)]
    } while (o < a && -1 == r);
    if (-1 == r)
      break;
    s += String.fromCharCode(t << 2 | (48 & r) >> 4);
    do {
      if (61 == (n = 255 & e.charCodeAt(o++)))
        return s;
      n = h[n]
    } while (o < a && -1 == n);
    if (-1 == n)
      break;
    s += String.fromCharCode((15 & r) << 4 | (60 & n) >> 2);
    do {
      if (61 == (i = 255 & e.charCodeAt(o++)))
        return s;
      i = h[i]
    } while (o < a && -1 == i);
    if (-1 == i)
      break;
    s += String.fromCharCode((3 & n) << 6 | i)
  }
  return s
}(a.data.info);

for (var c in a.data.encrypt_table) {
  var f = a.data.encrypt_table[c];
  if ("q" == f || "k" == f) {
    var p = a.data.key_table.pop();
    a.data.info = s[a.data.encrypt_table[c]](a.data.info, p)
  } else
    a.data.info = s[a.data.encrypt_table[c]](a.data.info)
}
var v = "";
for (c = 0; c < a.data.info.length; c++)
  v += String.fromCharCode(a.data.info[c]);

// fs.writeFileSync(path.join(__dirname, "data.txt"), v)
let sm = new stream.PassThrough()
sm.write(v)
sm.end()
sm.pipe(process.stdout)

function r(e, t) {
  var r = "";
  if ("object" == typeof e)
    for (var n in e)
      r += String.fromCharCode(e[n]);
  e = r || e;
  for (var i, o, a = new Uint8Array(e.length), s = t.length, n = 0; n < e.length; n++)
    o = n % s,
    i = e[n],
    i = i.toString().charCodeAt(0),
    a[n] = i ^ t.charCodeAt(o);
  return a
}

function n(e) {
  var t = "";
  if ("object" == typeof e)
    for (var r in e)
      t += String.fromCharCode(e[r]);
  e = t || e;
  var n = new Uint8Array(e.length);
  for (r = 0; r < e.length; r++)
    n[r] = e[r].toString().charCodeAt(0);
  var i, o, r = 0;
  for (r = 0; r < n.length; r++)
    0 != (i = n[r] % 3) && r + i < n.length && (o = n[r + 1],
      n[r + 1] = n[r + i],
      n[r + i] = o,
      r = r + i + 1);
  return n
}

function i(e) {
  var t = "";
  if ("object" == typeof e)
    for (var r in e)
      t += String.fromCharCode(e[r]);
  e = t || e;
  var n = new Uint8Array(e.length);
  for (r = 0; r < e.length; r++)
    n[r] = e[r].toString().charCodeAt(0);
  var r = 0,
    i = 0,
    o = 0,
    a = 0;
  for (r = 0; r < n.length; r++)
    o = n[r] % 2,
    o && r++,
    a++;
  var s = new Uint8Array(a);
  for (r = 0; r < n.length; r++)
    o = n[r] % 2,
    s[i++] = o ? n[r++] : n[r];
  return s
}

function o(e, t) {
  var r = 0,
    n = 0,
    i = 0,
    o = 0,
    a = "";
  if ("object" == typeof e)
    for (var r in e)
      a += String.fromCharCode(e[r]);
  e = a || e;
  var s = new Uint8Array(e.length);

  for (r = 0; r < e.length; r++)
    s[r] = e[r].toString().charCodeAt(0);
  for (r = 0; r < e.length; r++)
    if (0 != (o = s[r] % 5) && 1 != o && r + o < s.length && (i = s[r + 1],
        n = r + 2,
        s[r + 1] = s[r + o],
        s[o + r] = i,
        (r = r + o + 1) - 2 > n))
      for (; n < r - 2; n++)
        s[n] = s[n] ^ t.charCodeAt(n % t.length);
  for (r = 0; r < e.length; r++)
    s[r] = s[r] ^ t.charCodeAt(r % t.length);
  return s
}

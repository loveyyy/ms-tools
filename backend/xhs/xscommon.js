let window = global;
const args = process.argv.slice(2)
//上代理
function setProxy (proxyObjs) {
    for (let i = 0; i < proxyObjs.length; i++) {
        const handler = `{
          get: function(target, property, receiver) {
          if (property!="Math" && property!="isNaN"){
             if (target[property] && typeof target[property] !="string" &&  Object.keys(target[property]).length>3){
              }else{
            }}
            return target[property];
          },
          set: function(target, property, value, receiver) {
            
            return Reflect.set(...arguments);
          }
        }`;
        eval(`try {
            ${proxyObjs[i]};
            ${proxyObjs[i]} = new Proxy(${proxyObjs[i]}, ${handler});
        } catch (e) {
            ${proxyObjs[i]} = {};
            ${proxyObjs[i]} = new Proxy(${proxyObjs[i]}, ${handler});
        }`);
    }
}
navigator = {}
navigator.userAgent = 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36'
documentElement = {}
document = {
    documentElement: documentElement,
    cookie: args[0],
}
window.RegExp = function RegExp () { };
window.document = document;
window.window = window;
canvas = {
    getContext: function () { }
}
document.createElement = function (name) {
    return canvas
}

var canvas = {
    toDataURL: function toDataURL () {
    },
    getContext: function getContext (x) {
    },
    getAttribute: function getAttribute () {
        return null
    },
};
localStorage = {
    getItem: function getItem (x) {
        return null
    },
    removeItem: function removeItem (x) {
    }
};
setProxy(['window', 'document', ' navigator', 'screen', 'localStorage', 'location', 'RegExp'])

for (var P = [], A = "ZmserbBoHQtNP+wOcza/LpngG8yJq42KWYj0DSfdikx3VT16IlUAFM97hECvuRX5", R = 0, I = A.length; R < I; ++R)
    P[R] = A[R];
var encrypt_crc32 = function crc32(e) {
    for (var r, i = [], a = 0; a < 256; a++) {
        r = a;
        for (var s = 0; s < 8; s++)
            r = 1 & r ? 0xedb88320 ^ r >>> 1 : r >>> 1;
        i[a] = r
    }
    for (var u = -1, c = 0; c < e.length; c++)
        u = u >>> 8 ^ i[255 & (u ^ e.charCodeAt(c))];
    return (-1 ^ u) >>> 0
};

var O = function(e) {
    for (var r, i, a = 256, s = []; a--; s[a] = r >>> 0)
        for (i = 8,
                 r = a; i--; )
            r = 1 & r ? r >>> 1 ^ 0xedb88320 : r >>> 1;
    return function(e) {
        if ("string" == typeof e) {
            for (var r = 0, i = -1; r < e.length; ++r)
                i = s[255 & i ^ e.charCodeAt(r)] ^ i >>> 8;
            return -1 ^ i ^ 0xedb88320
        }
        for (var r = 0, i = -1; r < e.length; ++r)
            i = s[255 & i ^ e[r]] ^ i >>> 8;
        return -1 ^ i ^ 0xedb88320
    }
}()

function tripletToBase64(e) {
    return P[e >> 18 & 63] + P[e >> 12 & 63] + P[e >> 6 & 63] + P[63 & e]
}

function encodeChunk(e, r, i) {
    for (var a, s = [], u = r; u < i; u += 3)
        a = (e[u] << 16 & 0xff0000) + (e[u + 1] << 8 & 65280) + (255 & e[u + 2]),
            s.push(tripletToBase64(a));
    return s.join("")
}
function encodeUtf8(e) {
    for (var r = encodeURIComponent(e), i = [], a = 0; a < r.length; a++) {
        var s = r.charAt(a);
        if ("%" === s) {
            var u = parseInt(r.charAt(a + 1) + r.charAt(a + 2), 16);
            i.push(u),
                a += 2
        } else
            i.push(s.charCodeAt(0))
    }
    return i
}
function b64Encode(e) {
    for (var r, i = e.length, a = i % 3, s = [], u = 0, c = i - a; u < c; u += 16383)
        s.push(encodeChunk(e, u, u + 16383 > c ? c : u + 16383));
    return 1 === a ? (r = e[i - 1],
        s.push(P[r >> 2] + P[r << 4 & 63] + "==")) : 2 === a && (r = (e[i - 2] << 8) + e[i - 1],
        s.push(P[r >> 10] + P[r >> 4 & 63] + P[r << 2 & 63] + "=")),
        s.join("")
}

function getSigCount(e) {
    var r =  66;
    return e && (r++)
}

m = JSON.parse(args[0])
var a = O("".concat(m.x6).concat(m.x7).concat(m.x8))
m.x9 = a
f = m.x6 && m.x7 || ""
m.x10 = args[1]
m.x11 = 'normal'
console.log(b64Encode(encodeUtf8(JSON.stringify(m))))
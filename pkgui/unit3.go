// 由res2go IDE插件自动生成，不要编辑。
package pkgui

import (
    "github.com/ying32/govcl/vcl"
)

type TForm3 struct {
    *vcl.TForm
    ListBoxUserList *vcl.TListBox
    Label1          *vcl.TLabel
    ButtonEnter     *vcl.TButton
    ButtonCancel    *vcl.TButton

    //::private::
    TForm3Fields
}

var Form3 *TForm3




// vcl.Application.CreateForm(&Form3)

func NewForm3(owner vcl.IComponent) (root *TForm3)  {
    vcl.CreateResForm(owner, &root)
    return
}

var form3Bytes = []byte("\x54\x50\x46\x30\x06\x54\x46\x6F\x72\x6D\x33\x05\x46\x6F\x72\x6D\x33\x04\x4C\x65\x66\x74\x03\xFC\x01\x06\x48\x65\x69\x67\x68\x74\x03\x21\x01\x03\x54\x6F\x70\x03\x5A\x01\x05\x57\x69\x64\x74\x68\x03\x53\x01\x07\x43\x61\x70\x74\x69\x6F\x6E\x06\x06\xE5\xB8\xB3\xE8\x99\x9F\x0C\x43\x6C\x69\x65\x6E\x74\x48\x65\x69\x67\x68\x74\x03\x21\x01\x0B\x43\x6C\x69\x65\x6E\x74\x57\x69\x64\x74\x68\x03\x53\x01\x0A\x4C\x43\x4C\x56\x65\x72\x73\x69\x6F\x6E\x06\x07\x32\x2E\x32\x2E\x34\x2E\x30\x00\x08\x54\x4C\x69\x73\x74\x42\x6F\x78\x0F\x4C\x69\x73\x74\x42\x6F\x78\x55\x73\x65\x72\x4C\x69\x73\x74\x04\x4C\x65\x66\x74\x02\x10\x06\x48\x65\x69\x67\x68\x74\x03\xD1\x00\x03\x54\x6F\x70\x02\x40\x05\x57\x69\x64\x74\x68\x03\x38\x01\x0A\x49\x74\x65\x6D\x48\x65\x69\x67\x68\x74\x02\x00\x0B\x4D\x75\x6C\x74\x69\x53\x65\x6C\x65\x63\x74\x09\x08\x54\x61\x62\x4F\x72\x64\x65\x72\x02\x00\x00\x00\x06\x54\x4C\x61\x62\x65\x6C\x06\x4C\x61\x62\x65\x6C\x31\x04\x4C\x65\x66\x74\x02\x10\x06\x48\x65\x69\x67\x68\x74\x02\x14\x03\x54\x6F\x70\x02\x10\x05\x57\x69\x64\x74\x68\x02\x67\x07\x43\x61\x70\x74\x69\x6F\x6E\x06\x11\x43\x74\x72\x6C\x2B\xE5\xB7\xA6\xE9\x8D\xB5\xE9\x81\xB8\xE6\x93\x87\x0C\x46\x6F\x6E\x74\x2E\x43\x68\x61\x72\x53\x65\x74\x07\x13\x43\x48\x49\x4E\x45\x53\x45\x42\x49\x47\x35\x5F\x43\x48\x41\x52\x53\x45\x54\x0B\x46\x6F\x6E\x74\x2E\x48\x65\x69\x67\x68\x74\x02\xF0\x0A\x46\x6F\x6E\x74\x2E\x50\x69\x74\x63\x68\x07\x0A\x66\x70\x56\x61\x72\x69\x61\x62\x6C\x65\x0C\x46\x6F\x6E\x74\x2E\x51\x75\x61\x6C\x69\x74\x79\x07\x07\x66\x71\x44\x72\x61\x66\x74\x0B\x50\x61\x72\x65\x6E\x74\x43\x6F\x6C\x6F\x72\x08\x0A\x50\x61\x72\x65\x6E\x74\x46\x6F\x6E\x74\x08\x00\x00\x07\x54\x42\x75\x74\x74\x6F\x6E\x0B\x42\x75\x74\x74\x6F\x6E\x45\x6E\x74\x65\x72\x04\x4C\x65\x66\x74\x03\x88\x00\x06\x48\x65\x69\x67\x68\x74\x02\x25\x03\x54\x6F\x70\x02\x0A\x05\x57\x69\x64\x74\x68\x02\x59\x07\x43\x61\x70\x74\x69\x6F\x6E\x06\x06\xE7\xA2\xBA\xE8\xAA\x8D\x08\x54\x61\x62\x4F\x72\x64\x65\x72\x02\x01\x00\x00\x07\x54\x42\x75\x74\x74\x6F\x6E\x0C\x42\x75\x74\x74\x6F\x6E\x43\x61\x6E\x63\x65\x6C\x04\x4C\x65\x66\x74\x03\xEF\x00\x06\x48\x65\x69\x67\x68\x74\x02\x25\x03\x54\x6F\x70\x02\x0A\x05\x57\x69\x64\x74\x68\x02\x59\x07\x43\x61\x70\x74\x69\x6F\x6E\x06\x06\xE5\x8F\x96\xE6\xB6\x88\x08\x54\x61\x62\x4F\x72\x64\x65\x72\x02\x02\x00\x00\x00")

// 注册Form资源  
var _ = vcl.RegisterFormResource(Form3, &form3Bytes)

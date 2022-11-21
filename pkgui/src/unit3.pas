unit Unit3;

{$mode ObjFPC}{$H+}

interface

uses
  Classes, SysUtils, Forms, Controls, Graphics, Dialogs, StdCtrls;

type

  { TForm3 }

  TForm3 = class(TForm)
    ButtonEnter: TButton;
    ButtonCancel: TButton;
    Label1: TLabel;
    ListBoxUserList: TListBox;
    procedure ButtonCancelClick(Sender: TObject);
    procedure ButtonEnterClick(Sender: TObject);
    procedure FormClose(Sender: TObject; var CloseAction: TCloseAction);
    procedure FormCloseQuery(Sender: TObject; var CanClose: Boolean);
    procedure FormCreate(Sender: TObject);
    procedure FormShow(Sender: TObject);
  private

  public

  end;

var
  Form3: TForm3;

implementation

{$R *.lfm}

{ TForm3 }

procedure TForm3.ButtonEnterClick(Sender: TObject);
begin

end;

procedure TForm3.FormClose(Sender: TObject; var CloseAction: TCloseAction);
begin

end;

procedure TForm3.FormCloseQuery(Sender: TObject; var CanClose: Boolean);
begin

end;

procedure TForm3.FormCreate(Sender: TObject);
begin

end;

procedure TForm3.FormShow(Sender: TObject);
begin

end;

procedure TForm3.ButtonCancelClick(Sender: TObject);
begin

end;

end.


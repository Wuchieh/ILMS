unit Unit1;

{$mode objfpc}{$H+}

interface

uses
  Classes, SysUtils, Forms, Controls, Graphics, Dialogs, StdCtrls, ComCtrls,
  EditBtn, ValEdit, PairSplitter, Buttons, ExtCtrls, ActnList, Menus;

type

  { TForm1 }

  TForm1 = class(TForm)
    Button1: TButton;
    ButtonDeleteWork: TButton;
    ButtonLogout: TButton;
    ButtonDebug: TButton;
    ButtonLogin: TButton;
    ButtonUpdata: TButton;
    ComboBoxUserList: TComboBox;
    EditUsername: TEdit;
    EditPassword: TEdit;
    LabelPassword: TLabel;
    LabelUsername: TLabel;
    MemoConsole: TMemo;
    MenuItem1: TMenuItem;
    MenuItem2: TMenuItem;
    PopupMenu1: TPopupMenu;
    TreeViewClass: TTreeView;
    procedure ButtonDebugClick(Sender: TObject);
    procedure ButtonDeleteWorkClick(Sender: TObject);
    procedure ButtonLoginClick(Sender: TObject);
    procedure ButtonLogoutClick(Sender: TObject);
    procedure ButtonUpdataClick(Sender: TObject);
    procedure ComboBoxUserListChange(Sender: TObject);
    procedure EditPasswordKeyPress(Sender: TObject; var Key: char);
    procedure EditUsernameKeyPress(Sender: TObject; var Key: char);
    procedure FormClose(Sender: TObject; var CloseAction: TCloseAction);
    procedure FormCloseQuery(Sender: TObject; var CanClose: Boolean);
    procedure FormConstrainedResize(Sender: TObject; var MinWidth, MinHeight,
      MaxWidth, MaxHeight: TConstraintSize);
    procedure FormCreate(Sender: TObject);
    procedure FormResize(Sender: TObject);
    procedure FormShow(Sender: TObject);
    procedure MenuItem1Click(Sender: TObject);
    procedure MenuItem2Click(Sender: TObject);
    procedure MenuItem3Click(Sender: TObject);
    procedure TreeViewClassClick(Sender: TObject);
    procedure TreeViewClassDblClick(Sender: TObject);
    procedure TreeViewClassExpanded(Sender: TObject; Node: TTreeNode);
    procedure TreeViewClassMouseDown(Sender: TObject; Button: TMouseButton;
      Shift: TShiftState; X, Y: Integer);
  private

  public

  end;

var
  Form1: TForm1;

implementation

{$R *.lfm}

{ TForm1 }

procedure TForm1.FormCreate(Sender: TObject);
begin

end;

procedure TForm1.FormResize(Sender: TObject);
begin

end;

procedure TForm1.FormShow(Sender: TObject);
begin

end;

procedure TForm1.MenuItem1Click(Sender: TObject);
begin

end;

procedure TForm1.MenuItem2Click(Sender: TObject);
begin

end;

procedure TForm1.MenuItem3Click(Sender: TObject);
begin

end;

procedure TForm1.TreeViewClassClick(Sender: TObject);
begin

end;

procedure TForm1.TreeViewClassDblClick(Sender: TObject);
begin

end;

procedure TForm1.TreeViewClassExpanded(Sender: TObject; Node: TTreeNode);
begin

end;

procedure TForm1.TreeViewClassMouseDown(Sender: TObject; Button: TMouseButton;
  Shift: TShiftState; X, Y: Integer);
begin

end;

procedure TForm1.ButtonDebugClick(Sender: TObject);
begin

end;

procedure TForm1.ButtonDeleteWorkClick(Sender: TObject);
begin

end;

procedure TForm1.ButtonLoginClick(Sender: TObject);
begin

end;

procedure TForm1.ButtonLogoutClick(Sender: TObject);
begin

end;

procedure TForm1.ButtonUpdataClick(Sender: TObject);
begin

end;

procedure TForm1.ComboBoxUserListChange(Sender: TObject);
begin

end;

procedure TForm1.EditPasswordKeyPress(Sender: TObject; var Key: char);
begin

end;

procedure TForm1.EditUsernameKeyPress(Sender: TObject; var Key: char);
begin

end;

procedure TForm1.FormClose(Sender: TObject; var CloseAction: TCloseAction);
begin

end;

procedure TForm1.FormCloseQuery(Sender: TObject; var CanClose: Boolean);
begin

end;

procedure TForm1.FormConstrainedResize(Sender: TObject; var MinWidth,
  MinHeight, MaxWidth, MaxHeight: TConstraintSize);
begin

end;

end.


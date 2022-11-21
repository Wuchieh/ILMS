unit Unit2;

{$mode ObjFPC}{$H+}

interface

uses
  Classes, SysUtils, Forms, Controls, Graphics, Dialogs, StdCtrls;

type

  { TForm2 }

  TForm2 = class(TForm)
    ButtonDebug: TButton;
    ButtonMultipleSendWorks: TButton;
    ButtonUpload: TButton;
    ButtonInsert: TButton;
    ButtonDelete: TButton;
    ButtonSend: TButton;
    ComboBoxFiles: TComboBox;
    EditTitle: TEdit;
    LabelTitle: TLabel;
    LabelFiles: TLabel;
    MemoContent: TMemo;
    OpenDialogFiles: TOpenDialog;
    procedure ButtonDebugClick(Sender: TObject);
    procedure ButtonDeleteClick(Sender: TObject);
    procedure ButtonInsertClick(Sender: TObject);
    procedure ButtonSendClick(Sender: TObject);
    procedure ButtonMultipleSendWorksClick(Sender: TObject);
    procedure ButtonUploadClick(Sender: TObject);
    procedure ComboBoxFilesClick(Sender: TObject);
    procedure ComboBoxFilesEnter(Sender: TObject);
    procedure FormClose(Sender: TObject; var CloseAction: TCloseAction);
    procedure FormCloseQuery(Sender: TObject; var CanClose: Boolean);
    procedure FormConstrainedResize(Sender: TObject; var MinWidth, MinHeight,
      MaxWidth, MaxHeight: TConstraintSize);
    procedure FormCreate(Sender: TObject);
    procedure FormResize(Sender: TObject);
    procedure FormShow(Sender: TObject);
  private

  public

  end;

var
  Form2: TForm2;

implementation

{$R *.lfm}

{ TForm2 }

procedure TForm2.ButtonSendClick(Sender: TObject);
begin

end;

procedure TForm2.ButtonMultipleSendWorksClick(Sender: TObject);
begin

end;

procedure TForm2.ButtonUploadClick(Sender: TObject);
begin

end;

procedure TForm2.ComboBoxFilesClick(Sender: TObject);
begin

end;

procedure TForm2.ComboBoxFilesEnter(Sender: TObject);
begin

end;

procedure TForm2.FormClose(Sender: TObject; var CloseAction: TCloseAction);
begin

end;

procedure TForm2.FormCloseQuery(Sender: TObject; var CanClose: Boolean);
begin

end;

procedure TForm2.FormConstrainedResize(Sender: TObject; var MinWidth,
  MinHeight, MaxWidth, MaxHeight: TConstraintSize);
begin

end;

procedure TForm2.FormCreate(Sender: TObject);
begin

end;

procedure TForm2.FormResize(Sender: TObject);
begin

end;

procedure TForm2.FormShow(Sender: TObject);
begin

end;

procedure TForm2.ButtonDebugClick(Sender: TObject);
begin

end;

procedure TForm2.ButtonDeleteClick(Sender: TObject);
begin

end;

procedure TForm2.ButtonInsertClick(Sender: TObject);
begin

end;

end.


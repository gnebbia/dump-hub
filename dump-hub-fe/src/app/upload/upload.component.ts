/*
The MIT License (MIT)
Copyright (c) 2021 Davide Pataracchia
Permission is hereby granted, free of charge, to any person
obtaining a copy of this software and associated documentation
files (the "Software"), to deal in the Software without
restriction, including without limitation the rights to use,
copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the
Software is furnished to do so, subject to the following
conditions:
The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
OTHER DEALINGS IN THE SOFTWARE.
*/

import { HttpErrorResponse } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { ApiService } from '../api.service';

@Component({
  selector: 'app-upload',
  templateUrl: './upload.component.html',
  styleUrls: ['./upload.component.css']
})
export class UploadComponent implements OnInit {
  uploadForm = new FormGroup({
    pattern: new FormControl('', Validators.required),
    file: new FormControl('', Validators.required),
    columns: new FormControl('', Validators.required),
  });

  patternForm = new FormGroup({
    separator: new FormControl('', Validators.required),
    commentChar: new FormControl('', Validators.required)
  });

  uploadStatus: number = 0;
  editPatternModal: boolean = false;
  fileContent: string[] = [];
  fileContentRaw: any;
  selectedFile: any;

  previewContent: string[] = [];
  previewTable: string[][] = [];
  previewTableMaxCols: number = 0;

  constructor(
    private apiService: ApiService
  ) { }

  ngOnInit(): void {
    this.patternForm.setValue({
      separator: ":",
      commentChar: "#"
    });
    this.patternString();

    this.uploadForm.controls.pattern.disable();
    this.uploadForm.get('file')?.setValue(null);
    this.uploadForm.get('columns')?.setValue([]);
    this.previewContent = ["Select a text file to enable preview..."]

    this.onPatternChange();
    this.onCommentChange();
  }

  public onSubmit(): void {
    const formData = new FormData();
    formData.append('file', this.uploadForm.get('file')?.value);
    formData.append('pattern', this.uploadForm.get('pattern')?.value);
    formData.append('columns', this.uploadForm.get('columns')?.value);
    this.uploadStatus = 1;

    this.apiService.upload(formData).subscribe(
      (_) => {
        this.uploadStatus = 2;
        window.scroll(0, 0);
      },
      (err: HttpErrorResponse) => {     
        this.uploadStatus = -1;
      }
    );
  }

  public onFileSelect(event: any): void {
    this.uploadForm.get('file')!.setValue(null);
    this.uploadForm.get('columns')?.setValue([]);
    this.previewContent = ["Loading file..."]

    if (event.target.files.length > 0) {
      const file = event.target.files[0];
      if (!file.type.startsWith("text/")) {
        this.previewContent = ["Invalid file type.", "Please upload a text file."]
        this.uploadForm.get('file')!.setValue(null);
        return
      }
      this.uploadForm.get('file')!.setValue(file);
      this.readFile();
    }
  }

  public patternString(): void {
    var separator = this.patternForm.get('separator')?.value;
    var commentChar = this.patternForm.get('commentChar')?.value;

    var value = `{${separator}}{${commentChar}}`;
    this.uploadForm.get('pattern')!.setValue(value);
  }

  public parsePreview(): void {
    if (this.uploadForm.get('file')?.value == null) {
      return;
    }
    var separator = this.patternForm.get('separator')?.value;

    this.previewTable = [];
    this.previewTableMaxCols = 0;

    /* Get max column number */
    this.previewContent.forEach(content => {
      var values = content.replace(' ', '').split(separator);
      if (values.length > this.previewTableMaxCols) {
        this.previewTableMaxCols = values.length;
      }
    })

    /* Get table values */
    this.previewContent.forEach(content => {
      let tableRow: string[] = [];

      /* Init tableRow */
      for (var i = 0; i < this.previewTableMaxCols; i++) {
        tableRow[i] = "N/A";
      }

      var values = content.replace(' ', '').split(separator);
      for (var j = 0; j < values.length; j++) {
        if (values[j].length > 1) {
          tableRow[j] = values[j];
        }
      }

      this.previewTable.push(tableRow);
    });

    this.uploadForm.get('columns')?.setValue([]);
  }

  public loadingModal(): boolean {
    return this.uploadStatus != 0;
  }

  public counter(i: number) {
    return new Array(i);
  }

  public isSelected(colNumber: number) {
    var selected: number[] = this.uploadForm.get('columns')?.value;
    if (selected.indexOf(colNumber) !== -1) {
      return true;
    }
    return false;
  }

  public toggleColumn(colNumber: number) {
    var selected: number[] = this.uploadForm.get('columns')?.value;
    if (this.isSelected(colNumber)) {
      const index: number = selected.indexOf(colNumber);
      if (index !== -1) {
        selected.splice(index, 1);
      }

      this.uploadForm.get('columns')?.setValue(selected);
      return;
    }

    selected.push(colNumber)
    this.uploadForm.get('columns')?.setValue(selected);
  }

  private readFile(): void {
    this.fileContentRaw = null;
    var file = this.uploadForm.get('file')!.value;
    if (file == null) {
      return;
    }

    var reader: FileReader = new FileReader();
    reader.readAsText(file.slice(0, 8192));
    reader.onloadend = () => {
      this.fileContentRaw = reader.result
      this.processPreview();
    };

    reader.onerror = () => {
      this.uploadForm.get('file')!.setValue(null);
      this.previewContent = ["Unable to read the input file."]
      return
    }
  }

  private processPreview(): void {
    this.fileContent = this.fileContentRaw.split(/[\r\n]+/g);

    this.previewContent = [];
    for (var i = 0; i < this.fileContent.length - 1; i++) {
      if (this.previewContent.length >= 20) {
        break;
      }
      var commenctChar = this.patternForm.get('commentChar')?.value;
      if (this.fileContent[i].replace(' ', '').charAt(0) == commenctChar) {
        continue;
      }
      this.previewContent.push(this.fileContent[i]);
    }

    this.parsePreview();
  }

  private onPatternChange(): void {
    this.uploadForm.get('pattern')!.valueChanges.subscribe(_ => {
      this.parsePreview();
    });
  }

  private onCommentChange(): void {
    this.patternForm.get('commentChar')!.valueChanges.subscribe(_ => {
      this.processPreview();
    });
  }
}

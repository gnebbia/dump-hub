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

import { Component, OnInit } from '@angular/core';
import { FormControl, FormGroup } from '@angular/forms';
import { ApiService } from '../api.service';

interface SearchResponse {
  results?: any[],
  tot?: number
}

interface PagConfig {
  currentPage: number,
  pageSize: number,
  total: number
}

@Component({
  selector: 'app-main',
  templateUrl: './main.component.html',
  styleUrls: ['./main.component.css']
})
export class MainComponent implements OnInit {
  searchForm = new FormGroup({
    query: new FormControl(),
  });
  loadingResult: boolean = false;
  searchError: boolean = false;
  results: any[] = [];
  pagConfig: PagConfig;

  constructor(
    private apiService: ApiService
  ) {
    this.pagConfig = {
      currentPage: 1,
      pageSize: 20,
      total: 0
    }
  }

  ngOnInit(): void {
    this.initPaginator();
    this.onQueryChange();
    this.search();
  }

  public search(): void {
    var query = this.searchForm.get("query")?.value;
    this.loadingResult = true;

    this.apiService.search(query, this.pagConfig.currentPage)
      .subscribe(
        (data: SearchResponse) => {
          this.results = [];          
          if (data.results && data.tot) {
            this.results = data.results;
            this.pagConfig.total = data.tot;
            window.scroll(0, 0);
          }
          this.loadingResult = false;
        },
        _ => {
          this.results = [];
          this.loadingResult = false;
          this.searchError = true;
        }
      )
  }

  public pageChange(newPage: number): void {
    this.pagConfig.currentPage = newPage;
    this.search()
  }

  private initPaginator(): void {
    this.pagConfig = {
      currentPage: 1,
      pageSize: 20,
      total: 0
    }
  }

  private onQueryChange(): void {
    this.searchForm.get('query')!.valueChanges.subscribe(_ => {
      this.initPaginator();
      this.search();
    });
  }
}

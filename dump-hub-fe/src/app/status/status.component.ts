import { Component, OnDestroy, OnInit } from '@angular/core';
import { ApiService } from '../api.service';

interface History {
  date: string,
  filename: string,
  checksum: string,
  status: number,
}

interface HistoryData {
  results?: History[],
  tot?: number
}

interface PagConfig {
  currentPage: number,
  pageSize: number,
  total: number
}

@Component({
  selector: 'app-status',
  templateUrl: './status.component.html',
  styleUrls: ['./status.component.css']
})
export class StatusComponent implements OnInit, OnDestroy {
  uploadHistory: History[] = [];
  loadingHistory: boolean = false;
  historyError: boolean = false;
  apiInterval: any;
  pagConfig: PagConfig;
  deleteModal: boolean = false;
  toDelete: string = "";
  errorMessage: string = "Unable to retrieve history data";

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
    this.getHistory();

    this.apiInterval = setInterval(() => { this.getHistory() }, 30 * 1000);
  }

  ngOnDestroy() {
    clearInterval(this.apiInterval);
  }

  public getHistory(): void {
    this.loadingHistory = true;
    this.apiService.getHistory(this.pagConfig.currentPage).subscribe(
      (data: HistoryData) => {
        this.uploadHistory = [];
        this.pagConfig.total = 0;

        if (data.results && data.tot) {
          this.uploadHistory = data.results;
          this.pagConfig.total = data.tot;
        }
        this.loadingHistory = false;
      },
      _ => {
        this.errorMessage = "Unable to retrieve history data";
        this.historyError = true;
        this.loadingHistory = false;
      }
    );
  }

  public onDeleteRequest(checkSum: string): void {
    this.toDelete = checkSum;
    this.deleteModal = true;
  }

  public onDelete(): void {
    this.deleteHistory(this.toDelete);
    this.deleteModal = false;
  }

  public onDeleteCancel(): void {
    this.toDelete = "";
    this.deleteModal = false;
  }

  public deleteHistory(checkSum: string): void {
    this.apiService.delete(checkSum).subscribe(
      _ => {
        this.toDelete = "";
        this.deleteModal = false;
        this.getHistory();
      },
      _ => {
        this.errorMessage = "Unable to delete entries";
      }
    );
  }

  public pageChange(newPage: number): void {
    this.pagConfig.currentPage = newPage;
    this.getHistory()
  }

  private initPaginator(): void {
    this.pagConfig = {
      currentPage: 1,
      pageSize: 20,
      total: 0
    }
  }
}

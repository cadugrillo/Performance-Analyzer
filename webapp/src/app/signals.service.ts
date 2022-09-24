import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class SignalsService {

  constructor(private httpClient: HttpClient) {}

  parseSignals(file: any) {
    return this.httpClient.post(environment.gateway + '/performance-analyzer/signals/parse', file);
  }
}

export class ParsedSignals {
  aggregation!: string
  signalIds!: string[]
}

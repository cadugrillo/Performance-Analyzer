import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class SignalsService {

  constructor(private httpClient: HttpClient) {}

  parseSignals(file: any) {
    return this.httpClient.post(environment.gateway + '/performance-analyzer/signals/parse', file);
  }

  endpointResponse(file: any) {
    return this.httpClient.post(environment.gateway + '/performance-analyzer/signals/endresponse', file);
  }

  analizeData() {
    return this.httpClient.get(environment.gateway + '/performance-analyzer/signals/analyzedata')
  }
}

export class ParsedSignals {
  //aggregation!: string
  signalIds!: string[]
}

export class EndpointResponse {
  signals!: Signal[]
}

class Signal {
  signalId!: string
  legacySignalId!: number
  name!: string
  unit!: string
  type!: string
  aggregationId!: string
  values!: Value[]
}

class Value {
  timestamp!: EpochTimeStamp
  value!: any
}

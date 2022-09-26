import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { ParsedSignals } from './signals.service';


@Injectable({
  providedIn: 'root'
})
export class IotCoreEndpointService {

  constructor(private httpClient: HttpClient) {}

  getBulkSignals(recordAmount: number, millis: number, token: string, parsedSignals: ParsedSignals) {

    let httpOptions = {
      headers: new HttpHeaders({
        "Accept": "application/json",
        "Content-Type": "application/json",
        "Authorization": "Bearer "+token
      })
    };

    let endpoint_url = "https://api.industrial.voith.io/api/iot-core-signals/v3/signals/bulk/"+recordAmount+"/"+millis+"?disableConversion=false&useAvailabilityTime=false&normalization=none&direction=before";

    return this.httpClient.post(endpoint_url, parsedSignals);

  }
}

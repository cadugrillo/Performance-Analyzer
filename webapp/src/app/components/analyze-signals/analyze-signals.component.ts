import { Component, OnInit, ViewChild } from '@angular/core';
import { SignalsService, ParsedSignals } from '../../services/signals.service';
import { IotCoreEndpointService } from '../../services/iot-core-endpoint.service';
import { saveAs } from "file-saver";

@Component({
  selector: 'app-analyze-signals',
  templateUrl: './analyze-signals.component.html',
  styleUrls: ['./analyze-signals.component.css']
})
export class AnalyzeSignalsComponent implements OnInit {

  @ViewChild('xlsx') xlsx: any;
  @ViewChild('json') json: any;
  parsedSignals!: ParsedSignals;
  recordAmount!: number;
  millis!: number;
  token!: string;

  constructor(private SignalsService: SignalsService,
              private IotCoreEndpointService: IotCoreEndpointService) { }

  ngOnInit(): void {
  }

  ///////COMMON FUNCTIONS///////////////////

  toString(parsedSignals: Object): string {
    return JSON.stringify(parsedSignals, null, 4);
  }


  ///////PARSE SIGNALS FUNCTIONS/////////////////////////////

  importSigToParse() {
    this.xlsx.nativeElement.click();
   }
  
   onSigToParseAdded() {
    const xlsx = this.xlsx.nativeElement.files[0];
    this.SignalsService.parseSignals(xlsx).subscribe((data) => {
      this.parsedSignals = data as ParsedSignals;
      //console.log(this.parsedSignals);
    });
   }

  exportParsedSignals() {
    return saveAs(new Blob([JSON.stringify(this.parsedSignals, null, 2)], { type: 'JSON' }), 'parsedSignals.json');
  }


  ///////ENDPOINT RESPONSE FUNCTIONS/////////////////////////////

  importEndpResponse() {
    this.json.nativeElement.click();
   }
  
   onEndpResponseAdded() {
    const json = this.json.nativeElement.files[0];
    this.SignalsService.endpointResponse(json).subscribe((data) => {
      //this.parsedSignals = data as ParsedSignals;
      console.log(data);
    });
   }

   exportEndpointResponse() {
    return saveAs(new Blob([JSON.stringify(this.parsedSignals, null, 2)], { type: 'JSON' }), 'EndpointResponse.json');
  }


  ///////QUERY ENDPOINT FUNCTIONS/////////////////////////////

  queryEndpoint() {
    this.IotCoreEndpointService.getBulkSignals(this.recordAmount, this.millis, this.token, this.parsedSignals).subscribe((data) => {
      console.log(data);
    });
  }



}

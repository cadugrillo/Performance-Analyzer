import { Component, OnInit, ViewChild } from '@angular/core';
import { SignalsService, ParsedSignals, EndpointResponse, AnalyzedData } from '../../services/signals.service';
import { MessagePopupComponent } from '../message-popup/message-popup.component';
import { WaitPopupComponent } from '../wait-popup/wait-popup.component';
import { MatDialog } from '@angular/material/dialog';
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
  endpointResponse!: EndpointResponse;
  analyzedData!: AnalyzedData;
  TsInterval!: number;
  millis!: number;
  token!: string;

  constructor(private SignalsService: SignalsService,
              public dialog: MatDialog) { }

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
    this.dialog.open(WaitPopupComponent, {});
    const xlsx = this.xlsx.nativeElement.files[0];
    this.SignalsService.parseSignals(xlsx).subscribe((data) => {
      this.parsedSignals = data as ParsedSignals;
      this.dialog.closeAll();
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
    this.dialog.open(WaitPopupComponent, {});
    const json = this.json.nativeElement.files[0];
    this.SignalsService.endpointResponse(json).subscribe((data) => {
      this.endpointResponse = data as EndpointResponse;
      this.dialog.closeAll();
    });
   }

   exportEndpointResponse() {
    return saveAs(new Blob([JSON.stringify(this.endpointResponse, null, 2)], { type: 'JSON' }), 'EndpointResponse.json');
  }

  ////////ANALYZE DATA FUNCTIONS/////////////////////////////

  analyzeData(nrec: number) {
    this.dialog.open(WaitPopupComponent, {});
    this.SignalsService.analizeData(nrec).subscribe((data) => {
      this.analyzedData = data as AnalyzedData;
      this.dialog.closeAll();
      this.dialog.open(MessagePopupComponent, {data: {title: "Analysis Finished", text: "Check Results Tab!"}});
    });
  }

  exportAnalyzedData() {
    return saveAs(new Blob([JSON.stringify(this.analyzedData, null, 2)], { type: 'JSON' }), 'AnalyzedData.json');
  }
}

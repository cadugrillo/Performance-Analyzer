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

  @ViewChild('file') file: any;
  parsedSignals!: ParsedSignals;
  recordAmount!: number;
  millis!: number;
  token!: string;

  constructor(private SignalsService: SignalsService,
              private IotCoreEndpointService: IotCoreEndpointService) { }

  ngOnInit(): void {
  }

  importConfig() {
    this.file.nativeElement.click();
   }
  
   onFilesAdded() {
    const xlsx = this.file.nativeElement.files[0];
    this.SignalsService.parseSignals(xlsx).subscribe((data) => {
      this.parsedSignals = data as ParsedSignals;
      //console.log(this.parsedSignals);
    });
   }

   toString(parsedSignals: Object): string {
    return JSON.stringify(parsedSignals, null, 4);
  }

  exportParsedSignals() {
    return saveAs(new Blob([JSON.stringify(this.parsedSignals, null, 2)], { type: 'JSON' }), 'parsedSignals.json');
  }

  queryEndpoint() {
    this.IotCoreEndpointService.getBulkSignals(this.recordAmount, this.millis, this.token, this.parsedSignals).subscribe((data) => {
      console.log(data);
    });
  }

  exportEndpointResponse() {
    return saveAs(new Blob([JSON.stringify(this.parsedSignals, null, 2)], { type: 'JSON' }), 'EndpointResponse.json');
  }

}

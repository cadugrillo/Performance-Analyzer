import { Component, OnInit, ViewChild } from '@angular/core';
import { SignalsService, ParsedSignals } from '../signals.service';
import { saveAs } from "file-saver";

@Component({
  selector: 'app-parse-signals',
  templateUrl: './parse-signals.component.html',
  styleUrls: ['./parse-signals.component.css']
})
export class ParseSignalsComponent implements OnInit {

  @ViewChild('file') file: any;
  parsedSignals!: ParsedSignals;

  constructor(private SignalsService: SignalsService) { }

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

  exportMessages() {
    return saveAs(new Blob([JSON.stringify(this.parsedSignals, null, 2)], { type: 'JSON' }), 'parsedSignals.json');
  }

}

import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

export interface Node {
  id: number;
  x: number;
  y: number;
  neighbors: number[];
}

export interface Edge {
  from: number;
  to: number;
}

export interface GraphResponse {
  nodes: Node[];
  edges: Edge[];
  node_count: number;
}

@Injectable({
  providedIn: 'root',
})
export class GraphService {
  private baseUrl = 'http://localhost:6969/api';
  constructor(private http: HttpClient) {}

  generateGraph(nodesCount: number = 20): Observable<GraphResponse> {
    return this.http.get<GraphResponse>(
      `${this.baseUrl}/graph/generate?nodes=${nodesCount}`,
    );
  }

  startRouteDiscovery(src: number, dst: number): Observable<any> {
    const params = new HttpParams()
      .set('src', src.toString())
      .set('dst', dst.toString());
    return this.http.get(`${this.baseUrl}/simulation/start`, { params });
  }
}

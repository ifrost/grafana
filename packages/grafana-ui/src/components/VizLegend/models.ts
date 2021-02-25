export type LegendPlacement = 'bottom' | 'right';
export enum LegendDisplayMode {
  Hidden = 'hidden',
  List = 'list',
  Table = 'table',
}
export interface VizLegendOptions {
  calcs: string[];
  displayMode: LegendDisplayMode;
  placement: LegendPlacement;
}
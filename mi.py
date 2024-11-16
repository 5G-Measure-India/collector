from mobile_insight.analyzer import MsgLogger
from mobile_insight.monitor import OnlineMonitor

if __name__ == "__main__":
    src = OnlineMonitor()
    src.set_serial_port("/dev/ttyUSB0")
    src.set_baudrate(9600)

    for log in [
        # "5G_NR_NAS_SM_Plain_OTA_Incoming_Msg",
        # "5G_NR_NAS_SM_Plain_OTA_Outgoing_Msg",
        # "5G_NR_NAS_MM5G_State",
        "5G_NR_RRC_OTA_Packet",
        # "5G_NR_PDCP_UL_Control_Pdu",
        # "5G_NR_L2_UL_TB",
        # "5G_NR_L2_UL_BSR",
        # "5G_NR_RLC_DL_Stats",
        # "5G_NR_MAC_UL_TB_Stats",
        # "5G_NR_MAC_UL_Physical_Channel_Schedule_Report",
        # "5G_NR_MAC_PDSCH_Stats",
        # "5G_NR_MAC_RACH_Trigger",
        # "5G_NR_LL1_FW_Serving_FTL",
        "5G_NR_ML1_Searcher_Measurement_Database_Update_Ext",
        # "5G_NR_ML1_Serving_Cell_Beam_Management",
    ]:
        src.enable_log(log)

    logger = MsgLogger()
    logger.set_decode_format(MsgLogger.JSON)
    logger.set_source(src)

    src.run()
